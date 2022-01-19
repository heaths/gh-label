package export

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

type exportOptions struct {
	path   string
	format string

	// test
	client *github.Client
	io     *iostreams.IOStreams
}

// Make available for testing.
var opts *exportOptions

func ExportCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts = &exportOptions{}
	cmd := &cobra.Command{
		Use:   "export <path>",
		Short: `Export labels from the repository to <path>, or stdout if <path> is "-".`,
		Example: heredoc.Doc(`
			$ gh label export ./labels.csv
			$ gh label export ./labels.json
			$ gh label export --format csv -
		`),
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.format != "" {
				if format, err := github.SupportedOutputFormat(opts.format); err != nil {
					return err
				} else {
					opts.format = format
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.path = args[0]
			if opts.path != "-" {
				if opts.format == "" {
					opts.format = path.Ext(opts.path)
				}
			} else if opts.format == "" {
				return fmt.Errorf(`--format is required when <path> is "-"`)
			}

			if format, err := github.SupportedOutputFormat(opts.format); err != nil {
				return fmt.Errorf("%q has unsupported format %q, expected %v", opts.path, opts.format, github.OutputFormats())
			} else {
				opts.format = format
			}

			return export(globalOpts, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.format, "format", "", "", fmt.Sprintf("Format of the input to parse. One of %v. The default is the file extension.", github.OutputFormats()))

	return cmd
}

func export(globalOpts *options.GlobalOptions, opts *exportOptions) error {
	if opts.client == nil {
		owner, repo := globalOpts.Repo()
		cli := &github.Cli{
			Owner: owner,
			Repo:  repo,
		}
		opts.client = github.New(cli)
	}

	if opts.io == nil {
		opts.io = iostreams.System()
	}

	labels, err := opts.client.ListLabels("")
	if err != nil {
		return fmt.Errorf("failed to list labels; error: %w", err)
	}

	var w io.Writer
	if opts.path == "-" {
		w = opts.io.Out
	} else {
		if file, err := os.Create(opts.path); err != nil {
			return fmt.Errorf("failed to create file %q; error: %w", opts.path, err)
		} else {
			w = file
			defer file.Close()
		}
	}

	return labels.Write(github.OutputFormat(opts.format), w)
}
