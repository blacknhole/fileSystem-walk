package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "FilterNoExtension", root: "testdata",
			cfg:      config{"", 0, true, false},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			cfg:      config{".log", 0, true, false},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionNoMatch", root: "testdata",
			cfg:      config{".gz", 0, true, false},
			expected: ""},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			cfg:      config{".log", 10, true, false},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			cfg:      config{".log", 20, true, false},
			expected: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		t.Run(tc.name, func(t *testing.T) {
			if err := run(tc.root, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}

			result := buf.String()
			if result != tc.expected {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, result)
			}
		})
	}
}

func createTempDir(t *testing.T,
	files map[string]int) (tempDir string, cleanup func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}
	for k, n := range files {
		for i := 0; i < n; i++ {
			fname := fmt.Sprintf("file%d%s", i+1, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 0, nNoDelete: 10,
			expected: ""},
		{name: "DeleteExtensionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "", nDelete: 10, nNoDelete: 0,
			expected: ""},
		{name: "DeleteExtensionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 5, nNoDelete: 5,
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buf, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buf.String()
			if res != tc.expected {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n",
					tc.nNoDelete, len(filesLeft))
			}
		})
	}
}
