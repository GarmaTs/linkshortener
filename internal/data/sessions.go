package data

import (
	"sync"
	"time"
)

type Session struct {
	Username string
	Expiry   time.Time
}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

type SessionModel struct {
	mu       sync.Mutex
	sessions map[string]Session
}

func (sm *SessionModel) Set(token string, username string, expiry time.Time) Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := Session{
		Username: username,
		Expiry:   expiry,
	}
	sm.sessions[token] = session

	return session
}

func (sm *SessionModel) Get(token string) (Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var userSession Session
	userSession, exists := sm.sessions[token]
	if !exists {
		return userSession, ErrUnauthorized
	}

	if userSession.isExpired() {
		delete(sm.sessions, token)
		return userSession, ErrUnauthorized
	}

	return userSession, nil
}
