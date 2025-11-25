package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJoinSlicesDedup(t *testing.T) {
	t.Parallel()
	existing := []string{"a", "b"}
	newItems := []string{"b", "c"}
	joined := joinSlices(newItems, existing)

	if len(joined) != 3 {
		t.Fatalf("expected 3 elements after join, got %d", len(joined))
	}
	if !sliceContains(joined, "a") || !sliceContains(joined, "b") || !sliceContains(joined, "c") {
		t.Fatalf("joined slice missing expected elements: %v", joined)
	}
}

func TestSliceContains(t *testing.T) {
	t.Parallel()
	s := []string{"x", "y", "z"}
	if !sliceContains(s, "y") {
		t.Fatalf("expected slice to contain y")
	}
	if sliceContains(s, "a") {
		t.Fatalf("did not expect slice to contain a")
	}
}

func TestGetWeeksInLastSixMon(t *testing.T) {
	t.Parallel()
	if weeks := GetWeeksInLastSixMon(15); weeks != 3 {
		t.Fatalf("expected 3 weeks for 15 days, got %d", weeks)
	}
}

func TestGetDefaultGitEmailFromOverride(t *testing.T) {
	tmp := t.TempDir()
	conf := filepath.Join(tmp, ".gitconfig")
	content := "[user]\n\temail = test@example.com\n"
	if err := os.WriteFile(conf, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	t.Setenv("VISUALIZEDGIT_GITCONFIG", conf)

	if got := GetDefaultGitEmail(); got != "test@example.com" {
		t.Fatalf("expected email from config, got %q", got)
	}
}
