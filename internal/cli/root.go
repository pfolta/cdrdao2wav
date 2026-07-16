package cli

import (
	"github.com/spf13/cobra"

	"github.com/pfolta/cdrdao2audio"
)

func NewRootCommand() *cobra.Command {
	metadata := cdrdao2audio.GetMetadata()

	cmd := &cobra.Command{
		Use:                metadata.Name,
		Short:              "Convert a cdrdao dump to individual audio tracks",
		DisableSuggestions: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(NewVersionCommand(metadata))

	return cmd
}
