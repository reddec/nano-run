package worker

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"nano-run/services/blob"
	"nano-run/services/blob/fsblob"
	"nano-run/services/meta"
	"nano-run/services/meta/micrometa"
	"nano-run/services/queue"
	"nano-run/services/queue/microqueue"
)

type (
	CompleteHandler func(ctx context.Context, requestID string, info *meta.Request)
	ProcessHandler  func(ctx context.Context, requestID, attemptID string, info *meta.Request)
)

const (
	defaultAttempts        = 3
	defaultInterval        = 3 * time.Second
	minimalFailedCode      = 500
	nsRequest         byte = 0x00
	nsAttempt         byte = 0x01
)

func Default(location string) (*Worker, error) {
	path := filepath.Join(location, "blobs")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	valid, err := regexp.Compile("^[a-zA-Z0-9-]+$")
	if err != nil {
		return nil, err
	}
	storage := fsblob.NewCheck(path, func(id string) bool {
		return valid.MatchString(id)
	})

	taskQueue, err := microqueue.NewMicroQueue(filepath.Join(location, "queue"))
	if err != nil {
		return nil, err
	}
	requeue, err := microqueue.NewMicroQueue(filepath.Join(location, "requeue"))
	if err != nil {
		return nil, err
	}
	metaStorage, err := micrometa.NewMetaStorage(filepath.Join(location, "meta"))
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		_ = requeue.Close()
		_ = taskQueue.Close()
		_ = metaStorage.Close()
	}

	wrk, err := New(taskQueue, requeue, storage, metaStorage)
	if err != nil {
		cleanup()
		return nil, err
	}
	wrk.cleanup = cleanup
	return wrk, nil
}

func New(tasks, requeue queue.Queue, blobs blob.Blob, meta meta.Meta) (*Worker, error) {
	wrk := &Worker{
		queue:       tasks,
		requeue:     requeue,
		blob:        blobs,
		meta:        meta,
		reloadMeta:  make(chan struct{}, 1),
		maxAttempts: defaultAttempts,
		interval:    defaultInterval,
		concurrency: runtime.NumCPU(),
		handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNoContent)
		}),
	}
	err := wrk.init()
	if err != nil {
		return nil, err
	}
	return wrk, nil
}

type Worker struct {
	queue   queue.Queue
	requeue queue.Queue
	blob    blob.Blob
	meta    meta.Meta
	handler http.Handler
	cleanup func()

	onDead      CompleteHandler
	onSuccess   CompleteHandler
	onProcess   ProcessHandler
	maxAttempts int
	concurrency int
	reloadMeta  chan struct{}
	interval    time.Duration
	sequence    uint64
}

func (mgr *Worker) init() error {
	return mgr.meta.Iterate(func(id string, record meta.Request) error {
		if _, v, err := decodeID(id); err == nil && v > mgr.sequence {
			mgr.sequence = v
		} else if err != nil {
			log.Println("found broken id:", id, "-", err)
		}
		if !record.Complete {
			log.Println("found incomplete job", id)
			return mgr.queue.Push([]byte(id))
		}
		return nil
	})
}

// Cleanup internal resource.
func (mgr *Worker) Close() {
	if fn := mgr.cleanup; fn != nil {
		fn()
	}
}

// Enqueue request to storage, save meta-info to meta storage and push id into processing queue. Generated ID
// always unique and returns only in case of successful enqueue.
func (mgr *Worker) Enqueue(req *http.Request) (string, error) {
	id, err := mgr.saveRequest(req)
	if err != nil {
		return "", err
	}
	log.Println("new request saved:", id)
	err = mgr.queue.Push([]byte(id))
	return id, err
}

// Complete request manually.
func (mgr *Worker) Complete(requestID string) error {
	err := mgr.meta.Complete(requestID)
	if err != nil {
		return err
	}
	select {
	case mgr.reloadMeta <- struct{}{}:
	default:
	}
	return nil
}

func (mgr *Worker) OnSuccess(handler CompleteHandler) *Worker {
	mgr.onSuccess = handler
	return mgr
}

func (mgr *Worker) OnDead(handler CompleteHandler) *Worker {
	mgr.onDead = handler
	return mgr
}

func (mgr *Worker) OnProcess(handler ProcessHandler) *Worker {
	mgr.onProcess = handler
	return mgr
}

func (mgr *Worker) Handler(handler http.Handler) *Worker {
	mgr.handler = handler
	return mgr
}

func (mgr *Worker) HandlerFunc(fn http.HandlerFunc) *Worker {
	mgr.handler = fn
	return mgr
}

// Attempts number of 500x requests.
func (mgr *Worker) Attempts(max int) *Worker {
	mgr.maxAttempts = max
	return mgr
}

// Interval between attempts.
func (mgr *Worker) Interval(duration time.Duration) *Worker {
	mgr.interval = duration
	return mgr
}

// Concurrency limit (number of parallel tasks). Does not affect already running worker.
// 0 means num CPU.
func (mgr *Worker) Concurrency(num int) *Worker {
	mgr.concurrency = num
	if num == 0 {
		mgr.concurrency = runtime.NumCPU()
	}
	return mgr
}

// Meta information about requests.
func (mgr *Worker) Meta() meta.Meta {
	return mgr.meta
}

// Blobs storage (for large objects).
func (mgr *Worker) Blobs() blob.Blob {
	return mgr.blob
}

func (mgr *Worker) Run(global context.Context) error {
	if mgr.interval < 0 {
		return fmt.Errorf("negative interval")
	}
	if mgr.maxAttempts < 0 {
		return fmt.Errorf("negative attempts")
	}
	if mgr.handler == nil {
		return fmt.Errorf("nil handler")
	}
	if mgr.concurrency <= 0 {
		return fmt.Errorf("invalid concurrency number")
	}
	ctx, cancel := context.WithCancel(global)
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < mgr.concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer cancel()
			err := mgr.runQueue(ctx)
			if err != nil {
				log.Println("worker", i, "stopped due to error:", err)
			} else {
				log.Println("worker", i, "stopped")
			}
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		err := mgr.runReQueue(ctx)
		if err != nil {
			log.Println("re-queue process stopped due to error:", err)
		} else {
			log.Println("re-queue process stopped")
		}
	}()
	wg.Wait()
	return ctx.Err()
}

// Retry processing.
func (mgr *Worker) Retry(ctx context.Context, requestID string) (string, error) {
	info, err := mgr.meta.Get(requestID)
	if err != nil {
		return "", err
	}
	req, err := mgr.restoreRequest(ctx, requestID, info)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	return mgr.Enqueue(req)
}

func (mgr *Worker) call(ctx context.Context, requestID string, info *meta.Request) error {
	// caller should ensure that request id is valid
	req, err := mgr.restoreRequest(ctx, requestID, info)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	attemptID := encodeID(nsAttempt, uint64(len(info.Attempts))+1)

	req.Header.Set("X-Correlation-Id", requestID)
	req.Header.Set("X-Attempt-Id", attemptID)
	req.Header.Set("X-Attempt", strconv.Itoa(len(info.Attempts)+1))

	var header meta.AttemptHeader

	err = mgr.blob.Push(attemptID, func(out io.Writer) error {
		res := openResponse(out)
		started := time.Now()
		mgr.handler.ServeHTTP(res, req)
		header = res.meta
		header.StartedAt = started
		return nil
	})
	if err != nil {
		return err
	}

	info, err = mgr.meta.AddAttempt(requestID, attemptID, header)
	if err != nil {
		return err
	}

	mgr.requestProcessed(ctx, requestID, attemptID, info)
	if header.Code >= minimalFailedCode {
		return fmt.Errorf("500 code returned: %d", header.Code)
	}
	return nil
}

func (mgr *Worker) runQueue(ctx context.Context) error {
	for {
		err := mgr.processQueueItem(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (mgr *Worker) runReQueue(ctx context.Context) error {
	for {
		err := mgr.processReQueueItem(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (mgr *Worker) processQueueItem(ctx context.Context) error {
	bid, err := mgr.queue.Get(ctx)
	if err != nil {
		return err
	}
	id := string(bid)
	log.Println("processing request", id)
	info, err := mgr.meta.Get(id)
	if err != nil {
		return fmt.Errorf("get request %s meta info: %w", id, err)
	}
	if info.Complete {
		log.Printf("request %s already complete", id)
		return nil
	}
	err = mgr.call(ctx, id, info)
	if err == nil {
		mgr.requestSuccess(ctx, id, info)
		return nil
	}
	return mgr.requeueItem(ctx, id, info)
}

func (mgr *Worker) processReQueueItem(ctx context.Context) error {
	var item requeueItem

	data, err := mgr.requeue.Get(ctx)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &item)

	if err != nil {
		return err
	}

	d := time.Since(item.At)
	if d < mgr.interval {
		var ok = false
		for !ok {
			info, err := mgr.meta.Get(item.ID)
			if err != nil {
				return fmt.Errorf("re-queue: get meta %s: %w", item.ID, err)
			}
			if info.Complete {
				log.Printf("re-queue: %s already complete", item.ID)
				return nil
			}
			select {
			case <-time.After(mgr.interval - d):
				ok = true
			case <-mgr.reloadMeta:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return mgr.queue.Push([]byte(item.ID))
}

func (mgr *Worker) requeueItem(ctx context.Context, id string, info *meta.Request) error {
	if len(info.Attempts) >= mgr.maxAttempts {
		mgr.requestDead(ctx, id, info)
		log.Println("maximum attempts reached for request", id)
		return nil
	}
	data, err := json.Marshal(requeueItem{
		At: time.Now(),
		ID: id,
	})
	if err != nil {
		return err
	}
	return mgr.requeue.Push(data)
}

func (mgr *Worker) saveRequest(req *http.Request) (string, error) {
	id := encodeID(nsRequest, atomic.AddUint64(&mgr.sequence, 1))
	err := mgr.blob.Push(id, func(out io.Writer) error {
		_, err := io.Copy(out, req.Body)
		return err
	})
	if err != nil {
		return "", err
	}
	return id, mgr.meta.CreateRequest(id, req.Header, req.URL.RequestURI(), req.Method)
}

func (mgr *Worker) requestDead(ctx context.Context, id string, info *meta.Request) {
	err := mgr.meta.Complete(id)
	if err != nil {
		log.Println("failed complete (dead) request:", err)
	}
	if handler := mgr.onDead; handler != nil {
		handler(ctx, id, info)
	}
	log.Println("request", id, "completely failed")
}

func (mgr *Worker) requestSuccess(ctx context.Context, id string, info *meta.Request) {
	err := mgr.meta.Complete(id)
	if err != nil {
		log.Println("failed complete (success) request:", err)
	}
	if handler := mgr.onSuccess; handler != nil {
		handler(ctx, id, info)
	}
	log.Println("request", id, "complete successfully")
}

func (mgr *Worker) requestProcessed(ctx context.Context, id string, attemptID string, info *meta.Request) {
	if handler := mgr.onProcess; handler != nil {
		handler(ctx, id, attemptID, info)
	}
	log.Println("request", id, "processed with attempt", attemptID)
}

func (mgr *Worker) restoreRequest(ctx context.Context, requestID string, info *meta.Request) (*http.Request, error) {
	f, err := mgr.blob.Get(requestID)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, info.Method, info.URI, f)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	for k, v := range info.Headers {
		req.Header[k] = v
	}
	return req, nil
}

func encodeID(nsID byte, id uint64) string {
	var data [9]byte
	data[0] = nsID
	binary.BigEndian.PutUint64(data[1:], id)
	return strings.ToUpper(hex.EncodeToString(data[:]))
}

func decodeID(val string) (byte, uint64, error) {
	const idLen = 1 + 8
	hx, err := hex.DecodeString(val)
	if err != nil {
		return 0, 0, err
	}
	if len(hx) != idLen {
		return 0, 0, errors.New("too short")
	}
	n := binary.BigEndian.Uint64(hx[1:])
	return hx[0], n, nil
}

type requeueItem struct {
	At time.Time
	ID string
}
