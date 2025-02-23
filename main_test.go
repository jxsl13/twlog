package main

import (
	"context"
	"io"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/jxsl13/twlog/internal/testutils"
)

func TestWhoSaidCommand(t *testing.T) {
	ctx := context.TODO()
	cmd := NewRootCmd(ctx)

	archiveFolder := testutils.FilePath("testdata/subdir")

	concurrency := runtime.NumCPU()
	out, err := testutils.Execute(
		cmd,
		"--concurrency",
		strconv.Itoa(concurrency),
		"--search-dir",
		archiveFolder,
		"--include-archive",
		"who",
		"said",
		"--ips-only",
		"--deduplicate",
		"te[il]egram",
	)
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
	data, err := io.ReadAll(out)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	result := strings.TrimSpace(string(data))
	if result == "" {
		t.Fatalf("expected some output, got nothing")
	}
}

func TestWhatSaidCommand(t *testing.T) {
	ctx := context.TODO()
	cmd := NewRootCmd(ctx)

	archiveFolder := testutils.FilePath("testdata")

	concurrency := runtime.NumCPU()
	out, err := testutils.Execute(
		cmd,
		"--concurrency",
		strconv.Itoa(concurrency),
		"--search-dir",
		archiveFolder,
		"--include-archive",
		"what",
		"said",
		"--ips-only",
		"--deduplicate",
		"OP",
	)
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
	data, err := io.ReadAll(out)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	result := strings.TrimSpace(string(data))
	if result == "" {
		t.Fatalf("expected some output, got nothing")
	}
}
