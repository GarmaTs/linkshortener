package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrUnauthorized   = errors.New("user unauthorized")
)

type Models struct {
	Users interface {
		Insert(user *User) error
		GetByName(name string) (*User, error)
		Update(user *User) error
	}
	Sessions interface {
		Set(token string, username string, expiry time.Time) Session
		Get(token string) (Session, error)
	}
}

func NewModels(db *sql.DB, sessionz map[string]Session) Models {
	return Models{
		Users:    UserModel{DB: db},
		Sessions: &SessionModel{sessions: sessionz},
	}
}
