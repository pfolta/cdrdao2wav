// Copyright (c) 2026 Peter Folta
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
