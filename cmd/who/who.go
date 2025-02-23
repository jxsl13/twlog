package who

import (
	"github.com/jxsl13/twlog/internal/sharedcontext"
	"github.com/spf13/cobra"
)

func NewWhoCommand(root *sharedcontext.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "who",
		Short: "who is the subcomand which allows to search for players, nicknames and their IPs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewSaidCommand(root))
	return cmd
}
