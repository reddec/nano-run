// +build linux darwin

package main

import (
	"os"
	"syscall"
)

var signals = []os.Signal{syscall.SIGTERM, os.Interrupt}
