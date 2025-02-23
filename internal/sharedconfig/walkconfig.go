package sharedconfig

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/jxsl13/twlog/fswalk"
)

type WalkConfig struct {
	SearchDir       string         `koanf:"search.dir" short:"d" description:"directory to search for files recursively"`
	FileRegex       string         `koanf:"file.regex" short:"f" description:"regex to match files in the search dir"`
	FileRegexp      *regexp.Regexp `koanf:"-"`
	ArchiveRegex    string         `koanf:"archive.regex" short:"a" description:"regex to match archive files in the search dir"`
	ArchiveRegexp   *regexp.Regexp `koanf:"-"`
	IncludeArchives bool           `koanf:"include.archive" short:"A" description:"search inside archive files"`
	Concurrency     int            `koanf:"concurrency" short:"t" description:"number of concurrent workers to use"`
}

func NewWalkConfig() WalkConfig {
	return WalkConfig{
		SearchDir:    ".",
		FileRegex:    `.*\.log$`,
		ArchiveRegex: `\.(7z|bz2|gz|tar|xz|zip|xz|zst|lz)$`,
		Concurrency:  max(1, runtime.NumCPU()),
	}
}

func (cfg *WalkConfig) ToFSWalkConfig() fswalk.WalkConfig {
	return fswalk.WalkConfig{
		SearchDir:       cfg.SearchDir,
		FileRegexp:      cfg.FileRegexp,
		ArchiveRegexp:   cfg.ArchiveRegexp,
		IncludeArchives: cfg.IncludeArchives,
		Concurrency:     cfg.Concurrency,
	}
}

func (cfg *WalkConfig) Validate() error {
	if cfg.SearchDir == "" {
		return errors.New("search dir is required")
	}

	fi, err := os.Stat(cfg.SearchDir)
	if err != nil {
		return fmt.Errorf("invalid search dir: %w", err)
	}
	if !fi.IsDir() {
		return errors.New("search dir is not a directory")
	}

	if cfg.FileRegex == "" {
		return errors.New("file regex is required")
	}

	re, err := regexp.Compile(cfg.FileRegex)
	if err != nil {
		return fmt.Errorf("invalid file regex: %w", err)
	}
	cfg.FileRegexp = re

	if cfg.IncludeArchives || cfg.ArchiveRegex != "" {
		re, err = regexp.Compile(cfg.ArchiveRegex)
		if err != nil {
			return fmt.Errorf("invalid archive regex: %w", err)
		}
		cfg.ArchiveRegexp = re
	}

	if cfg.Concurrency < 1 {
		return errors.New("concurrency must be greater than 0")
	}

	return nil
}
