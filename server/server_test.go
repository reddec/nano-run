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

func testServer(t *testing.T, cfg server.Config, units map[string]server.Unit) *server.Server {
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

	srv, err := cfg.Create(context.Background())
	if !assert.NoError(t, err) {
		srv.Close()
		t.Fatal("failed to create server")
	}
	return srv
}

func Test_create(t *testing.T) {
	srv := testServer(t, server.DefaultConfig(), map[string]server.Unit{
		"hello": {
			Command: "echo hello world",
		},
	})
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/hello/", bytes.NewBufferString("hello world"))
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
	assert.NotEmpty(t, res.Header().Get("X-Correlation-Id"))
	assert.Equal(t, "/hello/"+res.Header().Get("X-Correlation-Id"), res.Header().Get("Location"))
	requestID := res.Header().Get("X-Correlation-Id")

	infoURL := res.Header().Get("Location")
	req = httptest.NewRequest(http.MethodGet, infoURL, nil)
	res = httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, requestID, res.Header().Get("X-Correlation-Id"))
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
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
		if assert.Equal(t, http.StatusNotFound, res.Code) {
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
