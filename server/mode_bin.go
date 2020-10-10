package server

import (
	"context"
	"io/ioutil"
	"log"
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
	command         string
	workDir         string
	shell           string
	environment     []string
	timeout         time.Duration
	gracefulTimeout time.Duration
}

func (bh *binHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	marker := &markerResponse{res: writer}
	ctx := request.Context()
	cmd := exec.Command(bh.shell, "-c", bh.command) //nolint:gosec

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

	err := bh.run(ctx, cmd)

	if codeReset, ok := writer.(interface{ Status(status int) }); ok && err != nil {
		codeReset.Status(http.StatusBadGateway)
	}

	if err != nil {
		writer.Header().Set("X-Return-Code", strconv.Itoa(cmd.ProcessState.ExitCode()))
		writer.WriteHeader(http.StatusBadGateway)
	} else {
		writer.Header().Set("X-Return-Code", strconv.Itoa(cmd.ProcessState.ExitCode()))
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (bh *binHandler) run(global context.Context, cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}

	var (
		gracefulTimeout <-chan time.Time
		ctx             context.Context
	)

	var gracefulTimer *time.Ticker
	if bh.gracefulTimeout > 0 {
		gracefulTimer = time.NewTicker(bh.gracefulTimeout)
		defer gracefulTimer.Stop()
		gracefulTimeout = gracefulTimer.C
	}

	if bh.timeout > 0 {
		t, cancel := context.WithTimeout(global, bh.timeout)
		defer cancel()
		ctx = t
	} else {
		ctx = global
	}

	var process = make(chan error, 1)

	go func() {
		defer close(process)
		process <- cmd.Wait()
	}()

	for {
		select {
		case <-gracefulTimeout:
			if err := internal.IntSignal(cmd); err != nil {
				log.Println("failed send signal to process:", err)
			} else {
				log.Println("sent graceful shutdown to process")
			}
			gracefulTimer.Stop()
			gracefulTimeout = nil
		case <-ctx.Done():
			if err := internal.KillSignal(cmd); err != nil {
				log.Println("failed send kill to process:", err)
			} else {
				log.Println("sent kill to process")
			}
			return <-process
		case err := <-process:
			return err
		}
	}
}

func (bh *binHandler) cloneEnv() []string {
	var cp = make([]string, len(bh.environment))
	copy(cp, bh.environment)
	return cp
}
