package data

import (
	"errors"
)

var (
	ErrDoesNotExist = errors.New("does not exist")
)

type InMemoryUserModel struct {
	users map[string]User
}

func (m InMemoryUserModel) Insert(user *User) error {
	_, err := m.GetByName(user.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrDoesNotExist):
			m.users[user.Name] = *user
			return nil
		default:
			return err
		}
	}

	return nil
}

func (m InMemoryUserModel) GetByName(name string) (*User, error) {
	u, exists := m.users[name]
	if !exists {
		return nil, ErrDoesNotExist
	}

	return &u, nil
}

func (m InMemoryUserModel) Update(user *User) error {
	u, err := m.GetByName(user.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrDoesNotExist):
			return nil
		default:
			return err
		}
	}
	u.Email = user.Email
	u.Password.hash = user.Password.hash

	m.users[u.Name] = *user

	return nil
}
