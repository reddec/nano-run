package ui

import "sync"

const sessionCookie = "session-id"

type Session interface {
	Login() string
	Valid() bool
}

type SessionStorage interface {
	Save(sessionID string, session Session)
	Get(sessionID string) (Session, bool)
	Delete(sessionID string)
}

type memorySessions struct {
	store sync.Map
}

func (ms *memorySessions) Save(sessionID string, session Session) {
	ms.store.Store(sessionID, session)
}

func (ms *memorySessions) Get(sessionID string) (Session, bool) {
	v, ok := ms.store.Load(sessionID)
	if !ok {
		return nil, ok
	}
	s, ok := v.(Session)
	return s, ok
}

func (ms *memorySessions) Delete(sessionID string) {
	ms.store.Delete(sessionID)
}
