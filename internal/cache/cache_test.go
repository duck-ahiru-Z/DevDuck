package cache

import (
	"path/filepath"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

func TestStoreRoundTrip(t *testing.T) {
	store := New(filepath.Join(t.TempDir(), "cache.json"))
	want := model.Explanation{Summary: "summary", Hints: []string{"hint"}}
	if err := store.Put("key", want); err != nil {
		t.Fatal(err)
	}
	got, found, err := store.Get("key")
	if err != nil || !found {
		t.Fatalf("get: found=%v err=%v", found, err)
	}
	if got.Summary != want.Summary || len(got.Hints) != 1 || got.Hints[0] != want.Hints[0] {
		t.Fatalf("got %#v", got)
	}
}

func TestKeyChangesWithInputs(t *testing.T) {
	err := model.ErrorInfo{Source: "python", Kind: "ValueError", Message: "bad"}
	first := Key(err, teaching.Beginner, "")
	if first == Key(err, teaching.Advanced, "") || first == Key(err, teaching.Beginner, "context") {
		t.Fatal("cache key did not include all inputs")
	}
}

func TestKeyRedactsSecretsBeforeHashing(t *testing.T) {
	withSecret := model.ErrorInfo{Raw: "token=secret-value"}
	redacted := model.ErrorInfo{Raw: "[REDACTED]"}
	if Key(withSecret, teaching.Beginner, "") != Key(redacted, teaching.Beginner, "") {
		t.Fatal("cache key was based on an unredacted secret")
	}
}
