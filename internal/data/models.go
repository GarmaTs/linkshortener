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
		GetList(username string) ([]*Url, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Sessions: &SessionModel{
			sessions: make(map[string]Session),
		},
		Urls: UrlModel{DB: db},
	}
}

func FakeNewModels() Models {
	return Models{
		Sessions: &SessionModel{
			sessions: make(map[string]Session),
		},
		Users: InMemoryUserModel{
			users: make(map[string]User),
		},
		Urls: InMemUrlModel{
			urls: make(map[string]Url),
		},
	}
}
