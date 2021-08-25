package main

import (
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
		Short: "List labels for the repository matching optional 'label' prefix",
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
			labels(query: $label, orderBy: {field: NAME, direction: ASC}, first: 5, after: $endCursor) {
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
		// Format as JSON array elements to work around https://github.com/cli/cli/issues/1268
		"--template", `{{range .data.repository.labels.nodes}}{{printf "{\"name\":\"%s\",\"color\":\"%s\",\"description\":\"%s\"}," .name .color ""}}{{end}}`,
		"-F", fmt.Sprintf("owner=%s", owner),
		"-F", fmt.Sprintf("repo=%s", repo),
		"-F", fmt.Sprintf("label=%s", opts.label),
		"-f", fmt.Sprintf("query=%s", query),
	}

	stdout, _, err := gh(args...)
	if err != nil {
		return fmt.Errorf("failed to list labels; error: %w", err)
	}

	// Delete the trailing comma from the formatted buffer and wrap as a JSON array.
	buffer := stdout.Bytes()
	buffer = append(append([]byte("["), buffer[:len(buffer)-1]...), []byte("]")...)

	type label struct {
		Name        string
		Color       string
		Description string
	}

	type labels []label

	var resp labels
	if err = json.Unmarshal(buffer, &resp); err != nil {
		return fmt.Errorf("failed to read label list; error: %w, stdout: %s", err, stdout.String())
	}

	colorizer := func(color string) func(string) string {
		return func(s string) string {
			return cs.HexToRGB(color, s)
		}
	}

	printer := utils.NewTablePrinter(io)
	for _, label := range resp {
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
