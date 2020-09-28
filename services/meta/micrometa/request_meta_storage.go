package micrometa

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger"

	"nano-run/internal"
	"nano-run/services/meta"
)

func WrapMetaStorage(db *badger.DB) *MicroMeta {
	return &MicroMeta{db: db, wrapped: true}
}

func NewMetaStorage(location string) (*MicroMeta, error) {
	name := filepath.Base(location)
	db, err := badger.Open(badger.DefaultOptions(location).
		WithLogger(internal.NanoLogger(log.New(os.Stderr, "["+name+"] ", log.LstdFlags))))
	if err != nil {
		return nil, err
	}
	return &MicroMeta{
		db:      db,
		wrapped: false,
	}, nil
}

type MicroMeta struct {
	db      *badger.DB
	wrapped bool
}

func (rms *MicroMeta) Get(requestID string) (*meta.Request, error) {
	var ans meta.Request
	return &ans, rms.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(requestID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &ans)
		})
	})
}

func (rms *MicroMeta) CreateRequest(requestID string, headers http.Header, uri string, method string) error {
	var record = meta.Request{
		CreatedAt: time.Now(),
		Attempts:  make([]meta.Attempt, 0),
		Complete:  false,
		Headers:   headers,
		URI:       uri,
		Method:    method,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return rms.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(requestID), data)
	})
}

func (rms *MicroMeta) AddAttempt(requestID, attemptID string, header meta.AttemptHeader) (*meta.Request, error) {
	var ans meta.Request
	err := rms.updateRequest(requestID, func(record *meta.Request) error {
		record.Attempts = append(record.Attempts, meta.Attempt{
			ID:            attemptID,
			CreatedAt:     time.Now(),
			AttemptHeader: header,
		})
		ans = *record
		return nil
	})
	return &ans, err
}

func (rms *MicroMeta) Complete(requestID string) error {
	return rms.updateRequest(requestID, func(record *meta.Request) error {
		record.Complete = true
		record.CompleteAt = time.Now()
		return nil
	})
}

func (rms *MicroMeta) Close() error {
	if rms.wrapped {
		return nil
	}
	return rms.db.Close()
}

func (rms *MicroMeta) Iterate(handler func(id string, record meta.Request) error) error {
	return rms.db.View(func(txn *badger.Txn) error {
		cfg := badger.DefaultIteratorOptions
		cfg.Reverse = true
		iter := txn.NewIterator(cfg)
		iter.Rewind()
		defer iter.Close()
		for iter.Valid() {
			id := string(iter.Item().Key())
			var rec meta.Request
			err := iter.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &rec)
			})
			if err != nil {
				return err
			}
			err = handler(id, rec)
			if err != nil {
				return err
			}
			iter.Next()
		}
		return nil
	})
}

func (rms *MicroMeta) updateRequest(requestID string, tx func(record *meta.Request) error) error {
	return rms.db.Update(func(txn *badger.Txn) error {
		data, err := txn.Get([]byte(requestID))
		if err != nil {
			return err
		}
		var record meta.Request
		err = data.Value(func(val []byte) error {
			return json.Unmarshal(val, &record)
		})
		if err != nil {
			return err
		}

		err = tx(&record)
		if err != nil {
			return err
		}

		value, err := json.Marshal(record)
		if err != nil {
			return err
		}
		return txn.Set([]byte(requestID), value)
	})
}
