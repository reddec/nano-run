package meta

import (
	"net/http"
	"time"
)

type Meta interface {
	Get(requestID string) (*Request, error)
	CreateRequest(requestID string, headers http.Header, uri string, method string) error
	AddAttempt(requestID, attemptID string, header AttemptHeader) (*Request, error)
	Complete(requestID string) error
	Iterate(handler func(id string, record Request) error) error
}

type AttemptHeader struct {
	Code    int         `json:"code"`
	Headers http.Header `json:"headers"`
}

type Attempt struct {
	AttemptHeader
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Request struct {
	CreatedAt  time.Time   `json:"created_at"`
	CompleteAt time.Time   `json:"complete_at,omitempty"`
	Attempts   []Attempt   `json:"attempts"`
	Headers    http.Header `json:"headers"`
	URI        string      `json:"uri"`
	Method     string      `json:"method"`
	Complete   bool        `json:"complete"`
}
