package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi" //nolint:gosec
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	"nano-run/server/internal"
	"nano-run/worker"
)

type Unit struct {
	Interval      time.Duration     `yaml:"interval,omitempty"`    // interval between attempts
	Attempts      int               `yaml:"attempts,omitempty"`    // maximum number of attempts
	Workers       int               `yaml:"workers,omitempty"`     // concurrency level - number of parallel requests
	Mode          string            `yaml:"mode,omitempty"`        // execution mode: bin, cgi or proxy
	WorkDir       string            `yaml:"workdir,omitempty"`     // working directory for the worker. if empty - temporary one will generated automatically
	Command       string            `yaml:"command"`               // command in a shell to execute
	Timeout       time.Duration     `yaml:"timeout,omitempty"`     // maximum execution timeout (enabled only for bin mode and only if positive)
	Shell         string            `yaml:"shell,omitempty"`       // shell to execute command in bin mode (default - /bin/sh)
	Environment   map[string]string `yaml:"environment,omitempty"` // custom environment for executable (in addition to system)
	MaxRequest    int64             `yaml:"max_request,omitempty"` // optional maximum HTTP body size (enabled if positive)
	Authorization struct {
		JWT struct {
			Enable bool `yaml:"enable"` // enable JWT verification
			JWT    `yaml:",inline"`
		} `yaml:"jwt,omitempty"` // HMAC256 JWT verification with shared secret

		QueryToken struct {
			Enable     bool `yaml:"enable"` // enable query-based token access
			QueryToken `yaml:",inline"`
		} `yaml:"query_token,omitempty"` // plain API tokens in request query params

		HeaderToken struct {
			Enable      bool `yaml:"enable"` // enable header-based token access
			HeaderToken `yaml:",inline"`
		} `yaml:"header_token,omitempty"` // plain API tokens in request header

		Basic struct {
			Enable bool `yaml:"enable"` // enable basic verification
			Basic  `yaml:",inline"`
		} `yaml:"basic,omitempty"` // basic authorization
	} `yaml:"authorization,omitempty"`
	name string
}

const (
	defaultRequestSize = 1 * 1024 * 1024 // 1MB
	defaultAttempts    = 3
	defaultInterval    = 5 * time.Second
	defaultWorkers     = 1
	defaultShell       = "/bin/sh"
	defaultMode        = "bin"
	defaultCommand     = "echo hello world"
	defaultName        = "main"
)

func DefaultUnit() Unit {
	return Unit{
		Interval:   defaultInterval,
		Attempts:   defaultAttempts,
		Workers:    defaultWorkers,
		MaxRequest: defaultRequestSize,
		Shell:      defaultShell,
		Mode:       defaultMode,
		Command:    defaultCommand,
		name:       defaultName,
	}
}

func (cfg Unit) Validate() error {
	var checks []string
	if cfg.Interval < 0 {
		checks = append(checks, "negative interval")
	}
	if cfg.Attempts < 0 {
		checks = append(checks, "negative attempts")
	}
	if cfg.Workers < 0 {
		checks = append(checks, "negative workers amount")
	}
	if !(cfg.Mode == "bin" || cfg.Mode == "cgi" || cfg.Mode == "proxy") {
		checks = append(checks, "unknown mode "+cfg.Mode)
	}
	if len(checks) == 0 {
		return nil
	}
	return errors.New(strings.Join(checks, ", "))
}

func (cfg Unit) SaveFile(file string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0600)
}

func Units(configsDir string) ([]Unit, error) {
	var configs []Unit
	err := filepath.Walk(configsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			return nil
		}
		unitName := strings.ReplaceAll(strings.Trim(path[len(configsDir):strings.LastIndex(path, ".")], "/\\"), string(filepath.Separator), "-")
		cfg := DefaultUnit()
		cfg.name = unitName
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return err
		}
		configs = append(configs, cfg)
		return nil
	})
	return configs, err
}

func Workers(workdir string, configurations []Unit) ([]*worker.Worker, error) {
	var ans []*worker.Worker
	for _, cfg := range configurations {
		log.Println("validating", cfg.name)
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("configuration invalid for %s: %w", cfg.name, err)
		}
		if cfg.Workers == 0 {
			cfg.Workers = runtime.NumCPU()
		}
		wrk, err := cfg.worker(workdir)
		if err != nil {
			for _, w := range ans {
				w.Close()
			}
			return nil, err
		}
		ans = append(ans, wrk)
	}
	return ans, nil
}

func Handler(units []Unit, workers []*worker.Worker) http.Handler {
	router := mux.NewRouter()
	for i, unit := range units {
		prefix := "/" + unit.name + "/"
		subRouter := router.PathPrefix(prefix).Subrouter()
		subRouter.Use(unit.enableAuthorization())
		internal.Expose(subRouter, workers[i])
	}
	return router
}

func Run(global context.Context, workers []*worker.Worker) error {
	if len(workers) == 0 {
		<-global.Done()
		return global.Err()
	}

	ctx, cancel := context.WithCancel(global)
	defer cancel()
	var wg sync.WaitGroup

	for _, wrk := range workers {
		wg.Add(1)
		go func(wrk *worker.Worker) {
			err := wrk.Run(ctx)
			if err != nil {
				log.Println("failed:", err)
			}
			wg.Done()
		}(wrk)
	}

	wg.Wait()
	return ctx.Err()
}

func (cfg Unit) worker(root string) (*worker.Worker, error) {
	handler, err := cfg.handler()
	if err != nil {
		return nil, err
	}
	workdir := filepath.Join(root, cfg.name)
	wrk, err := worker.Default(workdir)
	if err != nil {
		return nil, err
	}
	wrk = wrk.Attempts(cfg.Attempts).Interval(cfg.Interval).Concurrency(cfg.Workers).Handler(handler)
	return wrk, nil
}

func (cfg Unit) handler() (http.Handler, error) {
	handler, err := cfg.createRunner()
	if err != nil {
		return nil, err
	}
	if cfg.MaxRequest > 0 {
		handler = limitRequest(cfg.MaxRequest, handler)
	}
	//TODO: add authorization
	return handler, nil
}

func (cfg Unit) createRunner() (http.Handler, error) {
	switch cfg.Mode {
	case "bin":
		return &binHandler{
			command:     cfg.Command,
			workDir:     cfg.WorkDir,
			shell:       cfg.Shell,
			timeout:     cfg.Timeout,
			environment: append(os.Environ(), makeEnvList(cfg.Environment)...),
		}, nil
	case "cgi":
		return &cgi.Handler{
			Path:   cfg.Shell,
			Dir:    cfg.WorkDir,
			Env:    append(os.Environ(), makeEnvList(cfg.Environment)...),
			Logger: log.New(os.Stderr, "[cgi] ", log.LstdFlags),
			Args:   []string{"-c", cfg.Command},
			Stderr: os.Stderr,
		}, nil
	case "proxy":
		// proxy to static URL
		u, err := url.Parse(cfg.Command)
		if err != nil {
			return nil, err
		}
		return httputil.NewSingleHostReverseProxy(u), nil
	default:
		return nil, fmt.Errorf("unknown mode %s", cfg.Mode)
	}
}

func makeEnvList(content map[string]string) []string {
	var ans = make([]string, 0, len(content))
	for k, v := range content {
		ans = append(ans, k+"="+v)
	}
	return ans
}
