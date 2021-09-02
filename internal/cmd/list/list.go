package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

type listOptions struct {
	label string

	// test
	client *github.Client
	io     *iostreams.IOStreams
}

func ListCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts := &listOptions{}
	cmd := &cobra.Command{
		Use:   "list [label]",
		Short: "List labels for the repository matching optional 'label' substring in the name or description",
		Example: heredoc.Doc(`
			$ gh label list
			$ gh label list service
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if len(args) > 0 {
				opts.label = args[0]
			}
			return list(globalOpts, opts)
		},
	}

	globalOpts.ConfigureCommand(cmd)
	return cmd
}

func list(globalOpts *options.GlobalOptions, opts *listOptions) error {

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

	labels, err := opts.client.ListLabels(opts.label)
	if err != nil {
		return fmt.Errorf("failed to list labels; error: %w", err)
	}

	io := opts.io
	cs := io.ColorScheme()

	colorizer := func(color string) func(string) string {
		return func(s string) string {
			return cs.HexToRGB(color, s)
		}
	}

	if io.IsStdoutTTY() {
		fmt.Fprintf(io.Out, "Showing %d labels\n\n", len(labels))
	}

	printer := utils.NewTablePrinter(io)
	for _, label := range labels {
		color := label.Color
		printer.AddField(label.Name, nil, colorizer(color))
		if printer.IsTTY() {
			color = "#" + color
		}
		printer.AddField(color, nil, nil)
		printer.AddField(label.Description, nil, cs.ColorFromString("gray"))
		printer.EndRow()
	}
	_ = printer.Render()

	return nil
}
