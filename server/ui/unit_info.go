package ui

import (
	"os"

	"nano-run/server"
	"nano-run/services/meta"
	"nano-run/worker"
)

type unitInfo struct {
	Unit        server.Unit
	Worker      *worker.Worker
	CronEntries []*server.CronEntry
}

type historyRecord struct {
	ID   string
	Meta meta.Request
}

func (info unitInfo) History(max int) ([]historyRecord, error) {
	var res []historyRecord

	err := info.Worker.Meta().Iterate(func(id string, record meta.Request) error {
		if len(res) >= max {
			return os.ErrClosed
		}
		res = append(res, historyRecord{
			ID:   id,
			Meta: record,
		})
		return nil
	})
	if err == os.ErrClosed || err == nil {
		return res, nil
	}
	return nil, err
}
