package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/cli/safeexec"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "label",
	}
	rootCmd.AddCommand(listCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type listOptions struct {
	limit uint
}

func listCmd() *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels for the repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(opts)
		},
	}

	cmd.Flags().UintVarP(&opts.limit, "limit", "L", 30, "Maximum number of labels to fetch (default 30)")

	return cmd
}

func list(opts listOptions) error {
	io := iostreams.System()
	cs := io.ColorScheme()

	query := `query ($owner: String!, $repo: String!, $limit: Int) {
		repository(name: $repo, owner: $owner) {
			labels(first: $limit, orderBy: {field: NAME, direction: ASC}) {
				nodes {
					name
					color
					description
				}
			}
		}
	}`

	args := []string{
		"api",
		"graphql",
		"-F", "owner=:owner",
		"-F", "repo=:repo",
		"-F", fmt.Sprintf("limit=%d", opts.limit),
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

	var resp response
	if err = json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return fmt.Errorf("failed to read label list; error: %w, stdout: %s", err, stdout.String())
	}

	colorizer := func(color string) func(string) string {
		return func(s string) string {
			return cs.HexToRGB(color, s)
		}
	}

	printer := utils.NewTablePrinter(io)
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

	_ = printer.Render()

	return nil
}

func gh(args ...string) (stdout, stderr bytes.Buffer, err error) {
	bin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("cannot find gh; is it installed? err: %w", err)
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh; error: %w, stderr: %s", err, stderr.String())
		return
	}

	return
}
