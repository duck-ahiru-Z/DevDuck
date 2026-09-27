package credential

import "testing"

type fakeStore struct{ values map[string]string }

func (f *fakeStore) Get(_, account string) (string, error) {
	value, ok := f.values[account]
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}
func (f *fakeStore) Set(_, account, secret string) error { f.values[account] = secret; return nil }
func (f *fakeStore) Delete(_, account string) error      { delete(f.values, account); return nil }

func TestFakeStoreCRUD(t *testing.T) {
	store := &fakeStore{values: map[string]string{}}
	if _, err := store.Get("DevDuck", "gemini"); err != ErrNotFound {
		t.Fatalf("get error = %v", err)
	}
	if err := store.Set("DevDuck", "gemini", "test-key"); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("DevDuck", "gemini")
	if err != nil || got != "test-key" {
		t.Fatalf("got %q, err=%v", got, err)
	}
	if err := store.Delete("DevDuck", "gemini"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get("DevDuck", "gemini"); err != ErrNotFound {
		t.Fatalf("get after delete = %v", err)
	}
}
