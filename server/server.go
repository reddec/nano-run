package server

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	"nano-run/worker"
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

func (cfg Config) Create(global context.Context) (*Server, error) {
	units, err := Units(cfg.ConfigDirectory)
	if err != nil {
		return nil, err
	}
	workers, err := Workers(cfg.WorkingDirectory, units)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(global)
	srv := &Server{
		Handler: Handler(units, workers),
		workers: workers,
		units:   units,
		done:    make(chan struct{}),
		cancel:  cancel,
	}
	go srv.run(ctx)
	return srv, nil
}

func (cfg Config) Run(global context.Context) error {
	ctx, cancel := context.WithCancel(global)
	defer cancel()

	srv, err := cfg.Create(global)
	if err != nil {
		return err
	}
	defer srv.Close()

	server := http.Server{
		Addr:    cfg.Bind,
		Handler: srv,
	}

	done := make(chan struct{})

	go func() {
		defer cancel()
		<-ctx.Done()
		t, c := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
		_ = server.Shutdown(t)
		c()
		close(done)
	}()

	if cfg.TLS.Enable {
		err = server.ListenAndServeTLS(cfg.TLS.Cert, cfg.TLS.Key)
	} else {
		err = server.ListenAndServe()
	}
	cancel()
	<-done
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

type Server struct {
	http.Handler
	workers []*worker.Worker
	units   []Unit
	cancel  func()
	done    chan struct{}
	err     error
}

func (srv *Server) Close() {
	for _, wrk := range srv.workers {
		wrk.Close()
	}
	srv.cancel()
	<-srv.done
}

func (srv *Server) Err() error {
	return srv.err
}

func (srv *Server) run(ctx context.Context) {
	err := Run(ctx, srv.workers)
	if err != nil {
		log.Println("workers stopped:", err)
	}
	srv.err = err
	close(srv.done)
}
