package server

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WorkingDirectory string        `yaml:"working_directory"`
	ConfigDirectory  string        `yaml:"config_directory"`
	Bind             string        `yaml:"bind"`
	GracefulShutdown time.Duration `yaml:"graceful_shutdown"`
	TLS              struct {
		Enable bool   `yaml:"enable"`
		Cert   string `yaml:"cert"`
		Key    string `yaml:"key"`
	} `yaml:"tls,omitempty"`
}

const (
	defaultGracefulShutdown = 5 * time.Second
	defaultBind             = "127.0.0.1:8989"
)

func DefaultConfig() Config {
	var cfg Config
	cfg.Bind = defaultBind
	cfg.WorkingDirectory = filepath.Join("run")
	cfg.ConfigDirectory = filepath.Join("conf.d")
	cfg.GracefulShutdown = defaultGracefulShutdown
	return cfg
}

func (cfg Config) CreateDirs() error {
	err := os.MkdirAll(cfg.WorkingDirectory, 0755)
	if err != nil {
		return err
	}
	return os.MkdirAll(cfg.ConfigDirectory, 0755)
}

func (cfg *Config) LoadFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	if !filepath.IsAbs(cfg.WorkingDirectory) {
		cfg.WorkingDirectory = filepath.Join(filepath.Dir(file), cfg.WorkingDirectory)
	}
	if !filepath.IsAbs(cfg.ConfigDirectory) {
		cfg.ConfigDirectory = filepath.Join(filepath.Dir(file), cfg.ConfigDirectory)
	}
	return nil
}

func (cfg Config) SaveFile(file string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0600)
}

func (cfg Config) Run(global context.Context) error {
	units, err := Units(cfg.ConfigDirectory)
	if err != nil {
		return err
	}
	workers, err := Workers(cfg.WorkingDirectory, units)
	if err != nil {
		return err
	}
	defer func() {
		for _, wrk := range workers {
			wrk.Close()
		}
	}()
	handler := Handler(units, workers)

	ctx, cancel := context.WithCancel(global)

	server := http.Server{
		Addr:    cfg.Bind,
		Handler: handler,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		<-ctx.Done()
		t, c := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
		_ = server.Shutdown(t)
		c()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		err := Run(ctx, workers)
		if err != nil {
			log.Println("workers stopped:", err)
		}
	}()
	if cfg.TLS.Enable {
		err = server.ListenAndServeTLS(cfg.TLS.Cert, cfg.TLS.Key)
	} else {
		err = server.ListenAndServe()
	}
	cancel()
	wg.Wait()
	return err
}

func limitRequest(maxSize int64, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		defer body.Close()
		if request.ContentLength > maxSize {
			http.Error(writer, "too big request", http.StatusBadRequest)
			return
		}

		limiter := io.LimitReader(request.Body, maxSize)
		request.Body = ioutil.NopCloser(limiter)
		handler.ServeHTTP(writer, request)
	})
}
