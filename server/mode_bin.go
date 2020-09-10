package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"nano-run/server/internal"
)

type markerResponse struct {
	dataSent bool
	res      http.ResponseWriter
}

func (m *markerResponse) Header() http.Header {
	return m.res.Header()
}

func (m *markerResponse) Write(bytes []byte) (int, error) {
	m.dataSent = true
	return m.res.Write(bytes)
}

func (m *markerResponse) WriteHeader(statusCode int) {
	m.dataSent = true
	m.res.WriteHeader(statusCode)
}

type binHandler struct {
	command     string
	workDir     string
	shell       string
	environment []string
	timeout     time.Duration
}

func (bh *binHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	marker := &markerResponse{res: writer}

	ctx := request.Context()
	if bh.timeout > 0 {
		c, cancel := context.WithTimeout(ctx, bh.timeout)
		defer cancel()
		ctx = c
	}

	cmd := exec.CommandContext(ctx, bh.shell, "-c", bh.command) //nolint:gosec

	if bh.workDir == "" {
		tmpDir, err := ioutil.TempDir("", "")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(tmpDir)
		cmd.Dir = tmpDir
	} else {
		cmd.Dir = bh.workDir
	}

	var env = bh.cloneEnv()
	for k, v := range request.Header {
		ke := strings.ToUpper(strings.Replace(k, "-", "_", -1))
		env = append(env, ke+"="+strings.Join(v, ","))
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = request.Body
	cmd.Stdout = marker
	cmd.Env = env
	internal.SetBinFlags(cmd)
	err := cmd.Run()

	if marker.dataSent {
		return
	}
	if err != nil {
		writer.Header().Set("X-Return-Code", strconv.Itoa(cmd.ProcessState.ExitCode()))
		writer.WriteHeader(http.StatusBadGateway)
	} else {
		writer.Header().Set("X-Return-Code", strconv.Itoa(cmd.ProcessState.ExitCode()))
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (bh *binHandler) cloneEnv() []string {
	var cp = make([]string, len(bh.environment))
	copy(cp, bh.environment)
	return cp
}
