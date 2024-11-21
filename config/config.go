package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

func NewConfig() Config {
	return Config{
		SearchDir:   ".",
		FileRegex:   ".*\\.log",
		Deduplicate: false,
		Output:      FormatText,
	}
}

type Config struct {
	PhraseRegex  string         `koanf:"phrase.regex" short:"p" description:"regex to search for that a player said"`
	PhraseRegexp *regexp.Regexp `koanf:"-"`
	SearchDir    string         `koanf:"search.dir" short:"d" description:"directory to search for files recursively"`
	FileRegex    string         `koanf:"file.regex" short:"f" description:"regex to match files in the search dir"`
	FileRegexp   *regexp.Regexp `koanf:"-"`
	Deduplicate  bool           `koanf:"deduplicate" short:"D" description:"deduplicate objects based on all fields"`
	Extended     bool           `koanf:"extended" short:"e" description:"add two additional fields, file and id to the output"`
	IPsOnly      bool           `koanf:"ips.only" short:"i" description:"only print IP addresses"`
	Output       string         `koanf:"output" short:"o" description:"output format, one of 'json' or 'text'"`
}

func (cfg *Config) Validate() error {
	if cfg.PhraseRegex == "" {
		return errors.New("regex is required")
	}

	re, err := regexp.Compile(cfg.PhraseRegex)
	if err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}
	cfg.PhraseRegexp = re

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

	re, err = regexp.Compile(cfg.FileRegex)
	if err != nil {
		return fmt.Errorf("invalid file regex: %w", err)
	}
	cfg.FileRegexp = re

	allowed := []string{FormatJSON, FormatText}
	lOutput := strings.ToLower(cfg.Output)
	if !isOneOf(lOutput, allowed...) {
		return fmt.Errorf("invalid output format %q: must be one of %v", cfg.Output, allowed)
	}
	cfg.Output = lOutput

	if cfg.Extended && cfg.IPsOnly {
		return errors.New("extended and ips only flags are mutually exclusive")
	}

	return nil
}

func isOneOf(s string, values ...string) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}
