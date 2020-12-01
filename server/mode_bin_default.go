//+build !linux

package server

import "os/exec"

func setUser(cmd *exec.Cmd, user string) error {
	return nil
}
