package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
)

var (
	version = "dev"
	commit  = "dev"
)

type Config struct {
	Run    runCmd    `command:"run" description:"run single unit"`
	Server serverCmd `command:"server" description:"manage server"`
}

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "server", "run")
	}
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	parser.LongDescription = "Async webhook processor with minimal system requirements.\n\n" +
		"Author: Baryshnikov Aleksandr <owner@reddec.net>\n" +
		"Source code: https://github.com/reddec/nano-run\n" +
		"License: Apache 2.0\n" +
		"Version: " + version + "\n" +
		"Revision: " + commit
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func SignalContext() context.Context {
	gctx, closer := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, signals...)
		for range c {
			closer()
			break
		}
	}()
	return gctx
}
