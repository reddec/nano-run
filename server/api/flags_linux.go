package api

import (
	"os/exec"
	"syscall"
)

func SetBinFlags(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
	cmd.SysProcAttr.Setpgid = true
}
