package runner

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"

	"nano-run/server"
	"nano-run/server/ui"
	"nano-run/worker"
)

type Config struct {
	UIDirectory      string           `yaml:"ui_directory"`
	WorkingDirectory string           `yaml:"working_directory"`
	ConfigDirectory  string           `yaml:"config_directory"`
	Bind             string           `yaml:"bind"`
	GracefulShutdown time.Duration    `yaml:"graceful_shutdown"`
	Auth             ui.Authorization `yaml:"auth,omitempty"`
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
	cfg.UIDirectory = filepath.Join("ui")
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
	units, err := server.Units(cfg.ConfigDirectory)
	if err != nil {
		return nil, err
	}
	workers, err := server.Workers(cfg.WorkingDirectory, units)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(global)
	router := gin.Default()
	router.SetFuncMap(sprig.HtmlFuncMap())
	uiPath := filepath.Join(cfg.UIDirectory, "*.html")
	if v, err := filepath.Glob(uiPath); err == nil && len(v) > 0 {
		// little hack to prevent panic on empty UI
		router.LoadHTMLGlob(uiPath)
	}
	router.GET("/", func(gctx *gin.Context) {
		gctx.Redirect(http.StatusTemporaryRedirect, "ui")
	})
	uiGroup := router.Group("/ui/")
	ui.Attach(uiGroup, units, workers, cfg.Auth)
	uiGroup.Static("/static/", filepath.Join(cfg.UIDirectory, "static"))
	server.Attach(router.Group("/api/"), units, workers)

	srv := &Server{
		Handler: router,
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

	httpServer := http.Server{
		Addr:    cfg.Bind,
		Handler: srv,
	}

	done := make(chan struct{})

	go func() {
		defer cancel()
		<-ctx.Done()
		t, c := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
		_ = httpServer.Shutdown(t)
		c()
		close(done)
	}()

	if cfg.TLS.Enable {
		err = httpServer.ListenAndServeTLS(cfg.TLS.Cert, cfg.TLS.Key)
	} else {
		err = httpServer.ListenAndServe()
	}
	cancel()
	<-done
	return err
}

type Server struct {
	http.Handler
	workers []*worker.Worker
	units   []server.Unit
	cancel  func()
	done    chan struct{}
	err     error
}

func (srv *Server) Units() []server.Unit { return srv.units }

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
	err := server.Run(ctx, srv.workers)
	if err != nil {
		log.Println("workers stopped:", err)
	}
	srv.err = err
	close(srv.done)
}
