package main

import (
	"os"
	"testing"
	"time"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			f := filterOut(tc.path, []string{tc.ext}, tc.minSize, time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC), info)
			if f != tc.expected {
				t.Errorf("Expected %t, got %t instead\n", tc.expected, f)
			}
		})
	}
}
