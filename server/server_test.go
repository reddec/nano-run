package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"nano-run/server"
	"nano-run/server/runner"
	"nano-run/services/meta"
)

var tmpDir string

func TestMain(main *testing.M) {
	var err error
	tmpDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	code := main.Run()
	_ = os.RemoveAll(tmpDir)
	os.Exit(code)
}

func testServer(t *testing.T, cfg runner.Config, units map[string]server.Unit) *runner.Server {
	sub, err := ioutil.TempDir(tmpDir, "")
	if !assert.NoError(t, err) {
		t.Fatal("failed to create temp dir", err)
	}
	cfg.ConfigDirectory = filepath.Join(sub, "config")
	cfg.WorkingDirectory = filepath.Join(sub, "data")
	err = cfg.CreateDirs()
	if !assert.NoError(t, err) {
		t.Fatal("failed to create dirs", err)
	}

	for name, unit := range units {
		err = unit.SaveFile(filepath.Join(cfg.ConfigDirectory, name+".yaml"))
		if !assert.NoError(t, err) {
			t.Fatal("failed to create unit", name, ":", err)
		}
	}

	srv, err := cfg.Create(context.Background(), cfg.DefaultWaitTime)
	if !assert.NoError(t, err) {
		srv.Close()
		t.Fatal("failed to create server")
	}
	return srv
}

func Test_create(t *testing.T) {
	srv := testServer(t, runner.DefaultConfig(), map[string]server.Unit{
		"hello": {
			Command: "echo -n hello world",
		},
	})
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/hello/", bytes.NewBufferString("hello world"))
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
	assert.NotEmpty(t, res.Header().Get("X-Correlation-Id"))
	assert.Equal(t, "/api/hello/"+res.Header().Get("X-Correlation-Id"), res.Header().Get("Location"))
	requestID := res.Header().Get("X-Correlation-Id")

	infoURL := res.Header().Get("Location")
	t.Log("Location:", infoURL)
	req = httptest.NewRequest(http.MethodGet, infoURL, nil)
	res = httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, requestID, res.Header().Get("X-Correlation-Id"))
	assert.Contains(t, res.Header().Get("Content-Type"), "application/json")
	var info meta.Request
	err := json.Unmarshal(res.Body.Bytes(), &info)
	assert.NoError(t, err)

	// wait for result
	var resultLocation string
	for {
		req = httptest.NewRequest(http.MethodGet, infoURL+"/completed", nil)
		res = httptest.NewRecorder()
		srv.ServeHTTP(res, req)
		if res.Code == http.StatusMovedPermanently {
			resultLocation = res.Header().Get("Location")
			break
		}
		if !assert.Equal(t, http.StatusTooEarly, res.Code) {
			return
		}
		time.Sleep(time.Second)
	}

	req = httptest.NewRequest(http.MethodGet, resultLocation, nil)
	res = httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "hello world", res.Body.String())
}

func Test_retryIfDataReturnedInBinMode(t *testing.T) {
	srv := testServer(t, runner.DefaultConfig(), map[string]server.Unit{
		"hello": {
			Command: "echo hello world; exit 1",
		},
	})
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/hello/", bytes.NewBufferString("hello world"))
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
	assert.NotEmpty(t, res.Header().Get("X-Correlation-Id"))
	assert.Equal(t, "/api/hello/"+res.Header().Get("X-Correlation-Id"), res.Header().Get("Location"))
	location := res.Header().Get("Location")

	// wait for first result
	for {
		req = httptest.NewRequest(http.MethodGet, location, nil)
		res = httptest.NewRecorder()
		srv.ServeHTTP(res, req)
		if !assert.Equal(t, http.StatusOK, res.Code) {
			return
		}
		var info meta.Request
		err := json.Unmarshal(res.Body.Bytes(), &info)
		assert.NoError(t, err)
		if len(info.Attempts) == 0 {
			time.Sleep(time.Second)
			continue
		}
		atp := info.Attempts[0]
		assert.Equal(t, http.StatusBadGateway, atp.Code)
		assert.Equal(t, "1", atp.Headers.Get("X-Return-Code"))
		break
	}

}

func TestCron(t *testing.T) {
	srv := testServer(t, runner.DefaultConfig(), map[string]server.Unit{
		"hello": {
			Command: "echo hello world",
			Cron: []server.CronSpec{
				{Spec: "@every 1s"},
			},
		},
	})
	defer srv.Close()
	time.Sleep(time.Second + 100*time.Millisecond)

	var first *meta.Request
	err := srv.Workers()[0].Meta().Iterate(func(id string, record meta.Request) error {
		first = &record
		return nil
	})

	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, first.Complete)
	assert.Len(t, first.Attempts, 1)
}

func TestSync(t *testing.T) {
	srv := testServer(t, runner.DefaultConfig(), map[string]server.Unit{
		"hello": {
			Command: "echo hello world",
		},
	})
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPut, "/api/hello/", bytes.NewBufferString("hello world"))
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
}
