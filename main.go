package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jxsl13/twlog/cmd/who"
	"github.com/jxsl13/twlog/internal/sharedcontext"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cmd := NewRootCmd(ctx)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewRootCmd(ctx context.Context) *cobra.Command {
	root := sharedcontext.NewRoot(ctx)

	cmd := cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "twlog is a utility for analyzing Teeworlds server logs",
	}

	cmd.PersistentPreRunE = root.PersistentPreRunE(&cmd)
	cmd.PersistentPostRunE = root.PersistentPostRunE(&cmd)

	cmd.AddCommand(who.NewWhoCommand(root))
	return &cmd
}
