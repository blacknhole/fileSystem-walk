package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path string, exts []string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	ret := false
	for _, ext := range exts {
		if ext != "" {
			ret = true
			if filepath.Ext(path) == ext {
				return false
			}
		}
	}
	return ret
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, logger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	logger.Println(path)
	return nil
}

func archiveFile(path, root, destDir string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return nil
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("%s.gz", filepath.Base(path))
	targPath := filepath.Join(destDir, relDir, fname)

	if err := os.MkdirAll(filepath.Dir(targPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	gzW := gzip.NewWriter(out)
	gzW.Name = filepath.Base(path)

	if _, err := io.Copy(gzW, in); err != nil {
		return err
	}

	if err := gzW.Close(); err != nil {
		return err
	}

	return out.Close()
}
