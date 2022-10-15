package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	exts    []string
	size    int64
	date    time.Time
	list    bool
	del     bool
	wLog    io.Writer
	archive string
}

type stringSlice []string

func (s *stringSlice) String() string {
	return "String slice"
}

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	var exts stringSlice

	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	flag.Var(&exts, "ext", "File extention to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	date := flag.String("date", "1970-01-01", "Earliest file modification date")
	logFile := flag.String("log", "", "File to log")
	archive := flag.String("archive", "", "Archive directory")

	flag.Parse()

	var (
		f   = os.Stderr
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	d, err := time.Parse(time.RFC3339[:10], *date)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c := config{
		exts:    exts,
		size:    *size,
		date:    d,
		list:    *list,
		del:     *del,
		wLog:    f,
		archive: *archive,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	logger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, cfg.exts, cfg.size, cfg.date, info) {
				return nil
			}

			if cfg.list {
				return listFile(path, out)
			}

			if cfg.archive != "" {
				if err := archiveFile(path, root, cfg.archive); err != nil {
					return err
				}
			}

			if cfg.del {
				return delFile(path, logger)
			}

			return listFile(path, out)
		})
}
