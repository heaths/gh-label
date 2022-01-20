package _import

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

type importOptions struct {
	path   string
	format string

	// test
	client *github.Client
	fs     fs.FS
	io     *iostreams.IOStreams
}

// Make available for testing.
var opts *importOptions

func ImportCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts = &importOptions{}
	cmd := &cobra.Command{
		Use:   "import <path>",
		Short: `Import labels into the repository from <path>, or stdin if <path> is "-".`,
		Example: heredoc.Doc(`
			$ gh label import ./labels.csv
			$ gh label import ./labels.json
			$ gh label import --format csv -
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

			return _import(globalOpts, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.format, "format", "", "", fmt.Sprintf("Format of the input to parse. One of %v. The default is the file extension.", github.OutputFormats()))

	return cmd
}

func _import(globalOpts *options.GlobalOptions, opts *importOptions) error {
	if opts.client == nil {
		owner, repo := globalOpts.Repo()
		cli := &github.Cli{
			Owner: owner,
			Repo:  repo,
		}
		opts.client = github.New(cli)
	}

	if opts.fs == nil {
		pwd, err := os.Getwd()
		if err != nil {
			pwd = "/"
		}
		opts.fs = os.DirFS(pwd)
	}

	if opts.io == nil {
		opts.io = iostreams.System()
	}

	var r io.Reader
	if opts.path == "-" {
		r = opts.io.In
	} else {
		if file, err := opts.fs.Open(opts.path); err != nil {
			return fmt.Errorf("failed to open file %q; error: %w", opts.path, err)
		} else {
			r = file
			defer file.Close()
		}
	}

	labels, err := github.ReadLabels(github.OutputFormat(opts.format), r)
	if err != nil {
		return fmt.Errorf("failed to read labels; error: %w", err)
	}

	if opts.io.IsStdoutTTY() {
		fmt.Fprintf(opts.io.Out, "Importing %d label(s) from %q\n\n", len(labels), opts.path)
	}

	// TODO: Write progress bar if TTY.

	successes := 0
	failures := 0

	for _, label := range labels {
		if _, err := opts.client.CreateOrUpdateLabel(label); err != nil {
			failures++
			fmt.Fprintf(opts.io.ErrOut, "Failed to import label %q\n", label.Name)
			continue
		}

		successes++
	}

	if opts.io.IsStdoutTTY() {
		if failures > 0 {
			fmt.Fprintf(opts.io.ErrOut, "\n")
		}

		fmt.Fprintf(opts.io.Out, "Successfully imported %d, failed to import %d label(s)\n", successes, failures)
	}

	if successes == 0 {
		return errors.New("failed to import all labels")
	}

	return nil
}
