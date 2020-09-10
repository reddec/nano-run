package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	runtime "runtime"
	"strconv"
	"strings"
	"time"

	"nano-run/server"
)

type runCmd struct {
	Directory   string        `long:"directory" short:"d" env:"DIRECTORY" description:"Data directory" default:"run"`
	Interval    time.Duration `long:"interval" short:"i" env:"INTERVAL" description:"Requeue interval" default:"3s"`
	Attempts    int           `long:"attempts" short:"a" env:"ATTEMPTS" description:"Max number of attempts" default:"5"`
	Concurrency int           `long:"concurrency" short:"c" env:"CONCURRENCY" description:"Number of parallel worker (0 - mean number of CPU)" default:"0"`
	Mode        string        `long:"mode" short:"m" env:"MODE" description:"Running mode" default:"bin" choice:"bin" choice:"cgi" choice:"proxy"`
	Bind        string        `long:"bind" short:"b" env:"BIND" description:"Binding address" default:"127.0.0.1:8989"`
	Args        struct {
		Executable string   `arg:"executable" description:"path to binary to invoke or url" required:"yes"`
		Args       []string `arg:"args" description:"executable args"`
	} `positional-args:"yes"`
}

func (cfg *runCmd) Execute([]string) error {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	srv := server.DefaultConfig()
	srv.Bind = cfg.Bind
	srv.WorkingDirectory = cfg.Directory
	srv.ConfigDirectory = tmpDir

	unit := server.DefaultUnit()

	var params []string
	params = append(params, strconv.Quote(cfg.Args.Executable))
	for _, arg := range cfg.Args.Args {
		params = append(params, strconv.Quote(arg))
	}
	unit.Command = strings.Join(params, " ")
	unit.WorkDir, _ = os.Getwd()
	unit.Attempts = cfg.Attempts
	unit.Interval = cfg.Interval
	unit.Workers = cfg.concurrency()
	unit.Mode = cfg.Mode

	err = unit.SaveFile(filepath.Join(tmpDir, "main.yaml"))
	if err != nil {
		return err
	}
	return srv.Run(SignalContext())
}

func (cfg runCmd) concurrency() int {
	if cfg.Concurrency <= 0 {
		return runtime.NumCPU()
	}
	return cfg.Concurrency
}
