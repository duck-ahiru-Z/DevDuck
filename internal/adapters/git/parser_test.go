package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	adapter := NewAdapter()
	for _, command := range []string{"git", "git.exe", `/usr/bin/git`} {
		if !adapter.Detect(command, nil) {
			t.Errorf("Detect(%q)=false", command)
		}
	}
	if adapter.Detect("gitk", nil) {
		t.Error("gitk should not match")
	}
}

func TestParseFixtures(t *testing.T) {
	cases := []struct{ name, kind string }{
		{"not_repository", "RepositoryError"}, {"pathspec", "RefError"}, {"auth_failed", "AuthenticationError"}, {"publickey", "AuthenticationError"}, {"dns", "NetworkError"}, {"timeout", "NetworkError"}, {"non_fast_forward", "RefError"}, {"merge_conflict", "MergeConflict"}, {"local_changes", "WorktreeError"}, {"branch_exists", "RefError"}, {"remote_exists", "RefError"}, {"refspec", "RefError"}, {"unrelated", "RefError"}, {"identity", "ConfigError"}, {"dubious", "RepositoryError"}, {"unknown", "GitError"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok || info.Kind != tc.kind || info.Source != "git" {
				t.Fatalf("got %#v ok=%v", info, ok)
			}
		})
	}
}

func TestNonFastForwardIncludesContext(t *testing.T) {
	data, _ := os.ReadFile(filepath.Join("testdata", "non_fast_forward.txt"))
	info, ok := Parse(string(data))
	if !ok || info.Kind != "RefError" || (!strings.Contains(info.Message, "non-fast-forward") && !strings.Contains(info.Message, "remote contains work")) {
		t.Fatalf("got %#v", info)
	}
}
