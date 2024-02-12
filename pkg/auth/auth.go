package auth

import (
	"crypto/rand"
	"errors"
)

type AuthCookie struct {
	UserID         string
	ExpirationTime int64
}

type AuthStore struct {
	Store map[string]AuthCookie
}

var authenticationStore AuthStore

func (store *AuthStore) Insert(client AuthCookie) error {
	if _, ok := store.Store[client.UserID]; !ok {
		return errors.New("Failed to insert Client into AuthStore")
	}

	return nil
}

func GenerateToken() ([]byte, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
