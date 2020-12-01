package server

import (
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func setUser(cmd *exec.Cmd, userName string) error {
	if userName == "" {
		return nil
	}
	info, err := user.Lookup(userName)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(info.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(info.Gid)
	if err != nil {
		return err
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}
	return nil
}
