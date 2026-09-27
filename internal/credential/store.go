package credential

import (
	"errors"

	keyring "github.com/zalando/go-keyring"
)

var ErrNotFound = errors.New("credential not found")

type Store interface {
	Get(service, account string) (string, error)
	Set(service, account, secret string) error
	Delete(service, account string) error
}

type KeyringStore struct{}

func NewKeyringStore() *KeyringStore { return &KeyringStore{} }

func (s *KeyringStore) Get(service, account string) (string, error) {
	secret, err := keyring.Get(service, account)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	return secret, err
}

func (s *KeyringStore) Set(service, account, secret string) error {
	return keyring.Set(service, account, secret)
}

func (s *KeyringStore) Delete(service, account string) error {
	err := keyring.Delete(service, account)
	if errors.Is(err, keyring.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
