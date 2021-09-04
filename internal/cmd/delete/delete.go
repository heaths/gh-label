package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	name string

	// test
	client *github.Client
	io     *iostreams.IOStreams
}

func DeleteCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts := &deleteOptions{}
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete the label <name> from the repository",
		Example: heredoc.Doc(`
			$ gh label delete p1
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]

			return delete(globalOpts, opts)
		},
	}

	return cmd
}

func delete(globalOpts *options.GlobalOptions, opts *deleteOptions) error {
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

	if err := opts.client.DeleteLabel(opts.name); err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	if opts.io.IsStdoutTTY() {
		fmt.Fprintf(opts.io.Out, "Deleted label '%s'\n", opts.name)
	}

	return nil
}
