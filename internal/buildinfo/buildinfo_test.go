package buildinfo

import "testing"

func TestStringDevDefault(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, BuildDate
	defer func() { Version, Commit, BuildDate = oldVersion, oldCommit, oldDate }()
	Version, Commit, BuildDate = "dev", "unknown", "unknown"
	if got := String(); got != "DevDuck dev" {
		t.Fatalf("got %q", got)
	}
}

func TestStringBuildInfo(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, BuildDate
	defer func() { Version, Commit, BuildDate = oldVersion, oldCommit, oldDate }()
	Version, Commit, BuildDate = "v0.1.0", "abc1234", "2026-09-27"
	want := "DevDuck v0.1.0\ncommit: abc1234\nbuilt: 2026-09-27"
	if got := String(); got != want {
		t.Fatalf("got %q", got)
	}
}
