package microqueue

import (
	"context"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/dgraph-io/badger"

	"nano-run/internal"
)

var ErrEmptyQueue = errors.New("empty queue")

func WrapMicroQueue(db *badger.DB) (*MicroQueue, error) {
	mc := &MicroQueue{db: db, wrapped: true, notify: make(chan struct{}, 1), close: make(chan struct{})}
	return mc, mc.db.DropAll()
}

func NewMicroQueue(location string) (*MicroQueue, error) {
	name := filepath.Base(location)
	db, err := badger.Open(badger.DefaultOptions(location).WithTruncate(true).WithLogger(internal.NanoLogger(log.New(os.Stderr, "["+name+"] ", log.LstdFlags))))
	if err != nil {
		return nil, err
	}
	return &MicroQueue{db: db, notify: make(chan struct{}, 1), close: make(chan struct{})}, nil
}

// Always fresh queue with offloading to fs if no readers.
// Optimized for multiple readers and multiple writers with number of items limited by FS
// and each value should fit to RAM.
type MicroQueue struct {
	db       *badger.DB
	sequence uint64
	wrapped  bool
	notify   chan struct{}
	close    chan struct{}
}

func (mq *MicroQueue) Push(payload []byte) error {
	id := atomic.AddUint64(&mq.sequence, 1)
	var key [8]byte
	binary.BigEndian.PutUint64(key[:], id)
	err := mq.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key[:], payload)
	})
	if err != nil {
		return err
	}
	mq.sendNotify()
	return nil
}

func (mq *MicroQueue) pop() ([]byte, error) {
	var ans []byte
	err := mq.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   1,
			Reverse:        false,
			AllVersions:    false,
		})
		defer it.Close()
		it.Rewind()

		if !it.Valid() {
			return ErrEmptyQueue
		}
		v, err := it.Item().ValueCopy(ans)
		if err != nil {
			return err
		}
		ans = v

		return txn.Delete(it.Item().Key())
	})
	if err != nil {
		return nil, err
	}
	mq.sendNotify()
	return ans, nil
}

// Get blocking for new item in a queue.
func (mq *MicroQueue) Get(ctx context.Context) ([]byte, error) {
	mq.sendNotify()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-mq.close:
			return nil, errors.New("queue closed")
		case <-mq.notify:
			v, err := mq.pop()
			if err == nil {
				return v, nil
			}
			if errors.Is(err, ErrEmptyQueue) {
				continue
			}
			return nil, err
		}
	}
}

func (mq *MicroQueue) Close() error {
	close(mq.close)
	if mq.wrapped {
		return nil
	}
	return mq.db.Close()
}

func (mq *MicroQueue) sendNotify() {
	select {
	case mq.notify <- struct{}{}:
	default:
	}
}
