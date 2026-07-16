package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pfolta/cdrdao2audio"
)

type versionOptions struct {
	short bool
}

func NewVersionCommand(metadata cdrdao2audio.Metadata) *cobra.Command {
	opts := versionOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runVersion(cmd.OutOrStdout(), metadata, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.short, "short", "s", false, "show version number only")

	return cmd
}

func runVersion(
	w io.Writer,
	metadata cdrdao2audio.Metadata,
	opts versionOptions,
) error {
	if opts.short {
		_, err := fmt.Fprintln(w, strings.TrimPrefix(metadata.Version, "v"))
		return err
	}

	_, err := fmt.Fprintf(
		w,
		"%s version %s-%s-%s (%s)\n\n%s\n",
		metadata.Name,
		metadata.Version,
		metadata.OS,
		metadata.Arch,
		metadata.BuildDate,
		metadata.License,
	)

	return err
}
