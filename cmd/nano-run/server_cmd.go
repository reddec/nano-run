package main

import (
	"log"
	"os"
	"path/filepath"

	"nano-run/server"
	"nano-run/server/runner"
)

type serverCmd struct {
	Run  serverRunCmd  `command:"run" description:"run server"`
	Init serverInitCmd `command:"init" description:"initialize server"`
}

type serverInitCmd struct {
	Directory  string `short:"d" long:"directory" env:"DIRECTORY" description:"Target directory" default:"server"`
	ConfigFile string `long:"config-file" env:"CONFIG_FILE" description:"Config file name" default:"server.yaml"`
	NoSample   bool   `long:"no-sample" env:"NO_SAMPLE" description:"Do not create same file"`
}

func (cmd *serverInitCmd) Execute([]string) error {
	err := os.MkdirAll(cmd.Directory, 0755)
	if err != nil {
		return err
	}
	cfg := runner.DefaultConfig()
	err = cfg.SaveFile(filepath.Join(cmd.Directory, cmd.ConfigFile))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(cmd.Directory, cfg.ConfigDirectory), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(cmd.Directory, cfg.WorkingDirectory), 0755)
	if err != nil {
		return err
	}
	if !cmd.NoSample {
		unit := server.DefaultUnit()
		err = unit.SaveFile(filepath.Join(cmd.Directory, cfg.ConfigDirectory, "sample.yaml"))
		if err != nil {
			return err
		}
	}
	return nil
}

type serverRunCmd struct {
	Fail   bool   `short:"f" long:"fail" env:"FAIL" description:"Fail if no config file"`
	Config string `short:"c" long:"config" env:"CONFIG" description:"Configuration file" default:"server.yaml"`
}

func (cmd *serverRunCmd) Execute([]string) error {
	cfg := runner.DefaultConfig()
	err := cfg.LoadFile(cmd.Config)
	if os.IsNotExist(err) && !cmd.Fail {
		log.Println("no config file found - using transient default configuration")
		cfg.ConfigDirectory = filepath.Join("run", "conf.d")
		cfg.WorkingDirectory = filepath.Join("run", "data")
		err := cfg.CreateDirs()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	log.Println("configuration loaded")
	return cfg.Run(SignalContext())
}
