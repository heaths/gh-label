package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

type listOptions struct {
	label string
}

func listCmd(rootOpts *rootOptions) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list [label]",
		Short: "List labels for the repository matching optional 'label' substring in the name or description",
		Example: heredoc.Doc(`
			$ gh label list
			$ gh label list service
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.label = args[0]
			}
			return list(rootOpts, opts)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Parent().PersistentPreRunE != nil {
				return cmd.Parent().PersistentPreRunE(cmd, args)
			}
			return nil
		},
	}

	return cmd
}

func list(rootOpts *rootOptions, opts listOptions) error {
	io := iostreams.System()
	cs := io.ColorScheme()

	query := `query ($owner: String!, $repo: String!, $label: String, $endCursor: String) {
		repository(name: $repo, owner: $owner) {
			labels(query: $label, orderBy: {field: NAME, direction: ASC}, first: 100, after: $endCursor) {
				nodes {
					name
					color
					description
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	}`

	owner, repo := rootOpts.repoOverride()

	args := []string{
		"api",
		"graphql",
		"--paginate",
		"-F", fmt.Sprintf("owner=%s", owner),
		"-F", fmt.Sprintf("repo=%s", repo),
		"-F", fmt.Sprintf("label=%s", opts.label),
		"-f", fmt.Sprintf("query=%s", query),
	}

	stdout, _, err := gh(args...)
	if err != nil {
		return fmt.Errorf("failed to list labels; error: %w", err)
	}

	type response struct {
		Data struct {
			Repository struct {
				Labels struct {
					Nodes []struct {
						Name        string
						Color       string
						Description string
					}
				}
			}
		}
	}

	colorizer := func(color string) func(string) string {
		return func(s string) string {
			return cs.HexToRGB(color, s)
		}
	}

	printer := utils.NewTablePrinter(io)

	// Work around https://github.com/cli/cli/issues/1268 by splitting responses after cursor info.
	for _, data := range bytes.SplitAfter(stdout.Bytes(), []byte("}}}}}")) {
		if len(data) == 0 {
			break
		}

		var resp response
		if err = json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to read labels; error: %w, data: %s", err, data)
		}

		for _, label := range resp.Data.Repository.Labels.Nodes {
			color := label.Color
			printer.AddField(label.Name, nil, colorizer(color))
			if printer.IsTTY() {
				color = "#" + color
			}
			printer.AddField(color, nil, nil)
			printer.AddField(label.Description, nil, cs.ColorFromString("gray"))
			printer.EndRow()
		}

	}

	_ = printer.Render()

	return nil
}
