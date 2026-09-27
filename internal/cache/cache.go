package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

type Store struct {
	path string
}

type entry struct {
	Key         string            `json:"key"`
	Explanation model.Explanation `json:"explanation"`
}

func New(path string) *Store { return &Store{path: path} }

func Key(err model.ErrorInfo, level teaching.Level, contextText string) string {
	payload := struct {
		Error   model.ErrorInfo
		Level   teaching.Level
		Context string
	}{err, level, contextText}
	data, _ := json.Marshal(payload)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (s *Store) Get(key string) (model.Explanation, bool, error) {
	entries, err := s.read()
	if err != nil {
		return model.Explanation{}, false, err
	}
	for _, item := range entries {
		if item.Key == key {
			return item.Explanation, true, nil
		}
	}
	return model.Explanation{}, false, nil
}

func (s *Store) Put(key string, explanation model.Explanation) error {
	entries, err := s.read()
	if err != nil {
		return err
	}
	for i := range entries {
		if entries[i].Key == key {
			entries[i].Explanation = explanation
			return s.write(entries)
		}
	}
	return s.write(append(entries, entry{Key: key, Explanation: explanation}))
}

func (s *Store) read() ([]entry, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *Store) write(entries []entry) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
