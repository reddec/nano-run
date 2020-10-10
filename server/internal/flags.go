// +build !linux

package internal

import (
	"os"
	"os/exec"
)

func SetBinFlags(cmd *exec.Cmd) {

}

func IntSignal(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Interrupt)
}

func KillSignal(cmd *exec.Cmd) error {
	return cmd.Process.Signal(os.Kill)
}
