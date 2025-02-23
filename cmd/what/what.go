package what

import (
	"github.com/jxsl13/twlog/internal/sharedcontext"
	"github.com/spf13/cobra"
)

func NewWhatCommand(root *sharedcontext.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "what",
		Short: "what is the subcomand which allows to search what players did",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewSaidCommand(root))
	return cmd
}
