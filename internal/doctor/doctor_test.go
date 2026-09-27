package doctor

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type fakeStore struct{ configured bool }

func (f fakeStore) Get(_, _ string) (string, error) {
	if f.configured {
		return "secret", nil
	}
	return "", errors.New("not configured")
}
func (f fakeStore) Set(_, _, _ string) error { return nil }
func (f fakeStore) Delete(_, _ string) error { return nil }

func TestRunDoesNotPrintCredential(t *testing.T) {
	var out bytes.Buffer
	checker := Checker{Credential: fakeStore{configured: true}, LookPath: func(string) (string, error) { return "/tool", nil }, ConfigLoad: func() error { return nil }}
	if code := checker.Run(&out); code != 0 {
		t.Fatalf("code=%d", code)
	}
	if strings.Contains(out.String(), "secret") {
		t.Fatal("credential was printed")
	}
	if !strings.Contains(out.String(), "Gemini credential configured") {
		t.Fatal(out.String())
	}
}
