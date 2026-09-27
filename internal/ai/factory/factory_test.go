package factory

import "testing"

type fakeStore struct {
	value string
	err   error
}

func (f fakeStore) Get(string, string) (string, error) { return f.value, f.err }
func (f fakeStore) Set(string, string, string) error   { return nil }
func (f fakeStore) Delete(string, string) error        { return nil }

func TestNewProviderFromEnvRequiresKey(t *testing.T) {
	t.Setenv("DEVDUCK_AI_PROVIDER", "gemini")
	t.Setenv("GEMINI_API_KEY", "")
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("expected missing API key error")
	}
}

func TestNewProviderFromEnvUsesCredentialStore(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("DEVDUCK_AI_PROVIDER", "gemini")
	provider, err := NewProviderFromEnv(fakeStore{value: "stored-key"})
	if err != nil || provider == nil {
		t.Fatalf("provider=%v err=%v", provider, err)
	}
}

func TestNewProviderFromEnvFallsBackToEnvironment(t *testing.T) {
	t.Setenv("DEVDUCK_AI_PROVIDER", "gemini")
	t.Setenv("GEMINI_API_KEY", "environment-key")
	provider, err := NewProviderFromEnv()
	if err != nil || provider == nil {
		t.Fatalf("provider=%v err=%v", provider, err)
	}
}

func TestNewProviderFromEnvRejectsUnknownProvider(t *testing.T) {
	t.Setenv("DEVDUCK_AI_PROVIDER", "unknown")
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("expected unsupported provider error")
	}
}
