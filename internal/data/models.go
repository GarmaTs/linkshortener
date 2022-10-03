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
		Remove(token string)
	}
	Urls interface {
		Insert(url *Url, userName, fullUrl string) error
		GetOne(url *Url, shortUrl string) error
	}
}

func NewModels(db *sql.DB, sessions map[string]Session) Models {
	return Models{
		Users:    UserModel{DB: db},
		Sessions: &SessionModel{sessions: sessions},
		Urls:     UrlModel{DB: db},
	}
}
