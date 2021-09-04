package edit

import (
	"fmt"
	"regexp"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
	"github.com/heaths/gh-label/internal/utils"
	"github.com/spf13/cobra"
)

type editOptions struct {
	name        string
	color       string
	description string
	newName     string

	// test
	client *github.Client
	io     *iostreams.IOStreams
}

func EditCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts := &editOptions{}
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit the label <name> in the repository",
		Example: heredoc.Doc(`
			$ gh label edit general --new-name feedback
			$ gh label edit feedback --color c046ff --description "User feedback"
		`),
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.color != "" {
				if color, err := utils.ValidateColor(opts.color); err != nil {
					return fmt.Errorf(`invalid flag "color": %s`, err)
				} else {
					// Set color without "#" prefix.
					opts.color = color
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]

			return edit(globalOpts, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.color, "color", "c", "", `The color of the label with or without "#" prefix.`)
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "Description of the label.")
	cmd.Flags().StringVarP(&opts.newName, "new-name", "", "", "Rename the label to the given new name.")

	return cmd
}

func edit(globalOpts *options.GlobalOptions, opts *editOptions) error {
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

	label := github.EditLabel{
		Label: github.Label{
			Name:        opts.name,
			Color:       opts.color,
			Description: opts.description,
		},
		NewName: opts.newName,
	}

	updated, err := opts.client.UpdateLabel(label)
	if err != nil {
		return fmt.Errorf("failed to create label; error: %w", err)
	}

	re := regexp.MustCompile("^https://api.([^/]+)/repos/(.*)$")
	matches := re.FindStringSubmatch(updated.Url)

	if opts.io.IsStdoutTTY() {
		if label.Name != updated.Name {
			fmt.Fprintf(opts.io.Out, "Renamed label '%s' to '%s'\n\n", label.Name, updated.Name)
		} else {
			fmt.Fprintf(opts.io.Out, "Updated label '%s'\n\n", updated.Name)
		}
	}

	if len(matches) == 3 {
		fmt.Fprintf(opts.io.Out, "https://%s/%s\n", matches[1], matches[2])
	} else {
		fmt.Fprintln(opts.io.Out, updated.Url)
	}

	return nil
}
