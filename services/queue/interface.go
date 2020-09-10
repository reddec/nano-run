package queue

import "context"

type Queue interface {
	Push(payload []byte) error
	Get(ctx context.Context) ([]byte, error)
}
