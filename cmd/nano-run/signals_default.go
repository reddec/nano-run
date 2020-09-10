// +build !darwin,!linux

package main

import (
	"os"
)

var signals = []os.Signal{os.Interrupt}
