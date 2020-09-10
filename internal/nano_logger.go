package internal

import (
	"log"

	"github.com/dgraph-io/badger"
)

func NanoLogger(wrap *log.Logger) badger.Logger {
	return &nanoLogger{logger: wrap}
}

type nanoLogger struct {
	logger *log.Logger
}

func (nl *nanoLogger) Errorf(s string, i ...interface{}) {
	nl.logger.Printf("[error] "+s, i...)
}

func (nl *nanoLogger) Warningf(s string, i ...interface{}) {
	nl.logger.Printf("[warn] "+s, i...)
}

func (nl *nanoLogger) Infof(s string, i ...interface{}) {
	nl.logger.Printf("[info] "+s, i...)
}

func (nl *nanoLogger) Debugf(s string, i ...interface{}) {
	nl.logger.Printf("[debug] "+s, i...)
}
