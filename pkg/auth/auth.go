package auth

import "errors"

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
