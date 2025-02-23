package testutils

import (
	"bytes"

	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func FilePath(relative string, up ...int) string {
	offset := 1
	if len(up) > 0 && up[0] > 0 {
		offset = up[0]
	}
	_, file, _, ok := runtime.Caller(offset)
	if !ok {
		panic("failed to get caller")
	}
	if filepath.IsAbs(relative) {
		panic(fmt.Sprintf("%s is an absolute file path, must be relative to the current go source file", relative))
	}
	abs := filepath.Join(filepath.Dir(file), relative)
	return abs
}

func Execute(cmd *cobra.Command, args ...string) (out *bytes.Buffer, err error) {
	b := bytes.NewBuffer(nil)
	cmd.SetOut(b)
	cmd.SetArgs(args)
	err = cmd.Execute()
	if err != nil {
		return nil, err
	}
	return b, nil
}
