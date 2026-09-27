package factory

import "testing"

func TestNewProviderFromEnvRequiresKey(t *testing.T) {
	t.Setenv("DEVDUCK_AI_PROVIDER", "gemini")
	t.Setenv("GEMINI_API_KEY", "")
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("expected missing API key error")
	}
}

func TestNewProviderFromEnvRejectsUnknownProvider(t *testing.T) {
	t.Setenv("DEVDUCK_AI_PROVIDER", "unknown")
	if _, err := NewProviderFromEnv(); err == nil {
		t.Fatal("expected unsupported provider error")
	}
}
