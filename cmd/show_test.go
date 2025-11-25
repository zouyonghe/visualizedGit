package cmd

import (
	"testing"
	"time"
)

func TestCountDaysSinceDateOutOfRange(t *testing.T) {
	t.Parallel()
	oldDate := time.Now().AddDate(0, -7, 0) // older than six months
	if got := countDaysSinceDate(oldDate); got != outOfRange {
		t.Fatalf("expected %d for out-of-range date, got %d", outOfRange, got)
	}
}

func TestBuildColsAlignsFirstWeek(t *testing.T) {
	t.Parallel()

	commits := make(map[int]int)
	for i := 0; i <= 10; i++ {
		commits[i] = i
	}

	offset := calcOffset()
	cols := buildCols(commits)

	if len(cols) == 0 {
		t.Fatalf("expected at least one column")
	}

	expectedCol0 := make([]int, 0, 7)
	for i := 0; i < 7-offset-1; i++ {
		expectedCol0 = append(expectedCol0, 0)
	}
	for i := 0; i < offset+1; i++ {
		expectedCol0 = append(expectedCol0, commits[i])
	}

	col0, ok := cols[0]
	if !ok {
		t.Fatalf("expected column 0 to exist")
	}
	if len(col0) != 7 {
		t.Fatalf("expected column 0 length 7, got %d", len(col0))
	}
	for i := range col0 {
		if col0[i] != expectedCol0[i] {
			t.Fatalf("column 0 mismatch at index %d: expected %d, got %d", i, expectedCol0[i], col0[i])
		}
	}

	remaining := len(commits) - offset - 1
	expectedCols := 1 + (remaining+6)/7 // first col + padded weeks
	if len(cols) != expectedCols {
		t.Fatalf("expected %d columns, got %d", expectedCols, len(cols))
	}

	for i, c := range cols {
		if len(c) != 7 {
			t.Fatalf("column %d length %d, want 7", i, len(c))
		}
	}
}
