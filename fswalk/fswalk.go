package fswalk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sync"

	"github.com/jxsl13/twlog/archive"
	"github.com/jxsl13/twlog/ctxutils"
)

type WalkConfig struct {
	SearchDir       string
	IncludeArchives bool
	ArchiveRegexp   *regexp.Regexp
	FileRegexp      *regexp.Regexp
	Concurrency     int
}

func Walk(ctx context.Context, cfg WalkConfig, do func(filePath string, file io.Reader) error) error {
	cfg.Concurrency = max(1, cfg.Concurrency)

	ctx, cancelCause := context.WithCancelCause(ctx)
	defer cancelCause(errors.New("walk default canceled"))

	files := make([]string, 0, 16)
	archives := make([]string, 0, 1)

	entryDir, err := filepath.Abs(cfg.SearchDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of root dir or file: %w", err)
	}

	// collect log file and archive paths
	err = filepath.WalkDir(entryDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		err = ctxutils.Done(ctx)
		if err != nil {
			return err
		}

		// skip non-files
		if !info.Type().IsRegular() {
			return nil
		}

		if cfg.IncludeArchives && cfg.ArchiveRegexp != nil && cfg.ArchiveRegexp.MatchString(path) {
			archives = append(archives, path)
			return nil
		}

		if cfg.FileRegexp == nil || !cfg.FileRegexp.MatchString(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	slices.Sort(files)
	slices.Sort(archives)

	wg := &sync.WaitGroup{}

	concurrency := make(chan struct{}, cfg.Concurrency)

	wg.Add(len(files))
	for _, file := range files {
		exec := func() {
			concurrency <- struct{}{}
			defer func() {
				<-concurrency
				wg.Done()
			}()

			err := func(filePath string) error {
				f, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %w", filePath, err)
				}
				defer f.Close()

				return do(filePath, f)
			}(file)
			if err != nil {
				cancelCause(fmt.Errorf("error while processinf file %s: %w", file, err))
				return
			}
		}

		if cfg.Concurrency > 1 {
			// only run in parallel if concurrency is greater than 1
			go exec()
		} else {
			exec()
			err = ctxutils.Done(ctx)
			if err != nil {
				return err
			}
		}
	}

	wg.Add(len(archives))
	for _, file := range archives {
		exec := func() {
			concurrency <- struct{}{}
			defer func() {
				<-concurrency
				wg.Done()
			}()

			err = archive.Walk(file, func(path string, info fs.FileInfo, r io.Reader, err error) error {
				if err != nil {
					return err
				}

				err = ctxutils.Done(ctx)
				if err != nil {
					return err
				}

				if !info.Mode().IsRegular() {
					// skip dirs & symlinks
					return nil
				}

				if cfg.FileRegexp != nil && !cfg.FileRegexp.MatchString(path) {
					return nil
				}

				filePath := fmt.Sprintf("%s@%s", file, path)
				return do(filePath, r)
			})
			if err != nil {
				if errors.Is(err, archive.ErrUnsupportedArchive) {
					log.Printf("skipping unsupported archive: %s", file)
					return
				}
				cancelCause(fmt.Errorf("failed to walk archive %s: %w", file, err))
			}
		}

		if cfg.Concurrency > 1 {
			// only run in parallel if concurrency is greater than 1
			go exec()
		} else {
			exec()
			err = ctxutils.Done(ctx)
			if err != nil {
				return err
			}
		}
	}
	wg.Wait()

	err = ctxutils.Done(ctx)
	if err != nil {
		return err
	}
	return nil
}
