package sharedcontext

import (
	"context"
	"errors"
	"log"

	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/jxsl13/twlog/internal/sharedconfig"
	"github.com/spf13/cobra"
)

func NewRoot(ctx context.Context) *Root {
	ctx, cancelCause := context.WithCancelCause(ctx)
	return &Root{
		Ctx:         ctx,
		CancelCause: cancelCause,
		Format:      sharedconfig.NewFormatConfig(),
		Walk:        sharedconfig.NewWalkConfig(),
	}
}

type Root struct {
	Ctx         context.Context
	CancelCause context.CancelCauseFunc
	Format      sharedconfig.FormatConfig
	Walk        sharedconfig.WalkConfig
}

func (cli *Root) PersistentPreRunE(cmd *cobra.Command) func(*cobra.Command, []string) error {
	formatParser := cliconfig.RegisterFlags(&cli.Format, true, cmd, cliconfig.WithoutConfigFile())
	walkParser := cliconfig.RegisterFlags(&cli.Walk, true, cmd)
	return func(cmd *cobra.Command, args []string) error {
		log.SetOutput(cmd.ErrOrStderr()) // redirect log output to stderr

		return errors.Join(
			formatParser(),
			walkParser(),
		)
	}
}

func (cli *Root) PersistentPostRunE(_ *cobra.Command) func(*cobra.Command, []string) error {
	// could register stuff here
	return func(cmd *cobra.Command, args []string) error {
		cli.CancelCause(context.Canceled) // cleanup only
		return nil
	}
}
