package worker

import "sync/atomic"

func newTracker(id string) *Tracker {
	return &Tracker{
		id:   id,
		done: make(chan struct{}),
	}
}

type Tracker struct {
	id        string
	done      chan struct{}
	success   bool
	attemptID string
	finished  int32
}

func (t *Tracker) ID() string {
	return t.id
}

func (t *Tracker) Done() <-chan struct{} {
	return t.done
}

func (t *Tracker) Success() bool {
	return t.success
}

func (t *Tracker) Attempt() string {
	return t.attemptID
}

func (t *Tracker) close() {
	if atomic.CompareAndSwapInt32(&t.finished, 0, 1) {
		close(t.done)
	}
}

func (t *Tracker) ok(attemptID string) {
	t.attemptID = attemptID
	t.success = true
	t.close()
}

func (t *Tracker) failed() {
	t.success = false
	t.close()
}
