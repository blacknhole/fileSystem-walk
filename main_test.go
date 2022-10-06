package main

import (
	"bytes"
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
			cfg:      config{"", 0, true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			cfg:      config{".log", 0, true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionNoMatch", root: "testdata",
			cfg:      config{".gz", 0, true},
			expected: ""},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			cfg:      config{".log", 10, true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			cfg:      config{".log", 20, true},
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
