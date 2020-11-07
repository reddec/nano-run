package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/robfig/cron/v3"

	"nano-run/worker"
)

type CronSpec struct {
	Spec        string            `yaml:"spec"`                   // cron tab spec with seconds precision
	Name        string            `yaml:"name,omitempty"`         // optional name to distinguish in logs and ui
	Headers     map[string]string `yaml:"headers,omitempty"`      // headers in simulated request
	Content     string            `yaml:"content,omitempty"`      // content in simulated request
	ContentFile string            `yaml:"content_file,omitempty"` // content file to read for request content
}

func (cs CronSpec) Validate() error {
	_, err := cron.ParseStandard(cs.Spec)
	return err
}

func (cs *CronSpec) Label(def string) string {
	if cs.Name != "" {
		return cs.Name
	}
	return def
}

func (cs *CronSpec) Request(requestPath string) (*http.Request, error) {
	var src io.ReadCloser
	if cs.Content != "" || cs.ContentFile == "" {
		src = ioutil.NopCloser(bytes.NewReader([]byte(cs.Content)))
	} else if f, err := os.Open(cs.ContentFile); err != nil {
		return nil, err
	} else {
		src = f
	}

	req, err := http.NewRequest(http.MethodPost, requestPath, src)
	if err != nil {
		_ = src.Close()
	}
	return req, err
}

type CronEntry struct {
	Name   string
	Spec   CronSpec
	Worker *worker.Worker
	Config Unit
	ID     cron.EntryID
}

func (ce *CronEntry) Unit() Unit { return ce.Config }

// Cron initializes cron engine and registers all required worker schedules to it.
func Cron(workers []*worker.Worker, configs []Unit) ([]*CronEntry, *cron.Cron, error) {
	engine := cron.New()
	var entries []*CronEntry
	for i, wrk := range workers {
		cfg := configs[i]
		for i, cronSpec := range cfg.Cron {
			name := cfg.Name() + "/" + cronSpec.Label(strconv.Itoa(i))
			entry := &CronEntry{
				Spec:   cronSpec,
				Worker: wrk,
				Config: cfg,
				Name:   name,
			}
			id, err := engine.AddJob(cronSpec.Spec, entry)
			if err != nil {
				return nil, nil, fmt.Errorf("cron record %s: %w", name, err)
			}
			entry.ID = id
			entries = append(entries, entry)
		}
	}
	return entries, engine, nil
}

func (ce *CronEntry) Run() {
	req, err := ce.Spec.Request(ce.Config.Path())
	if err != nil {
		log.Println("failed create cron", ce.Name, "request:", err)
		return
	}
	id, err := ce.Worker.Enqueue(req)
	if err != nil {
		log.Println("failed enqueue cron", ce.Name, "job:", err)
		return
	}
	log.Println("enqueued cron", ce.Name, "job with id", id)
}
