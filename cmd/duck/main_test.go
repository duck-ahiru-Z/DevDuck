package main

import "testing"

func TestAuthInvalidArgumentsReturnNonZero(t *testing.T) {
	if code := handleAuth([]string{"unknown"}); code == 0 {
		t.Fatal("expected invalid auth arguments to fail")
	}
}

func TestConfigInvalidArgumentsReturnNonZero(t *testing.T) {
	if code := handleConfig([]string{"unknown"}); code == 0 {
		t.Fatal("expected invalid config arguments to fail")
	}
	if code := handleConfig(nil); code == 0 {
		t.Fatal("expected missing config arguments to fail")
	}
}
