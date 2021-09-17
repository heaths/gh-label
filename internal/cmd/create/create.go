package create

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

type createOptions struct {
	name        string
	color       string
	description string

	// test
	client *github.Client
	io     *iostreams.IOStreams
}

func CreateCmd(globalOpts *options.GlobalOptions) *cobra.Command {
	opts := &createOptions{}
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create the label <name> in the repository",
		Example: heredoc.Doc(`
			$ gh label create feedback
			$ gh label create p1 --color e00808
			$ gh label create p2 --color "#ffa501" --description "Affects more than a few users"
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

			return create(globalOpts, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.color, "color", "c", "", `The color of the label with or without "#" prefix. A random color will be assigned if not specified.`)
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "Description of the label.")

	return cmd
}

func create(globalOpts *options.GlobalOptions, opts *createOptions) error {
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

	if opts.color == "" {
		opts.color = utils.RandomColor()
	}

	label := github.Label{
		Name:        opts.name,
		Color:       opts.color,
		Description: opts.description,
	}

	label, err := opts.client.CreateLabel(label)
	if err != nil {
		return fmt.Errorf("failed to create label; error: %w", err)
	}

	re := regexp.MustCompile("^https://api.([^/]+)/repos/(.*)$")
	matches := re.FindStringSubmatch(label.URL)

	if opts.io.IsStdoutTTY() {
		fmt.Fprintf(opts.io.Out, "Created label '%s'\n\n", label.Name)
	}

	if len(matches) == 3 {
		fmt.Fprintf(opts.io.Out, "https://%s/%s\n", matches[1], matches[2])
	} else {
		fmt.Fprintln(opts.io.Out, label.URL)
	}

	return nil
}
