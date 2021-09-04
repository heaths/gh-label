package github

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/cli/safeexec"
)

type Cli struct {
	Owner string
	Repo  string
}

func (cli *Cli) CreateLabel(label Label) (bytes.Buffer, error) {
	args := []string{
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
		"-F", fmt.Sprintf("name=%s", label.Name),
		"-f", fmt.Sprintf("color=%s", label.Color),
		"-F", fmt.Sprintf("description=%s", label.Description),
		"/repos/:owner/:repo/labels",
	}

	stdout, _, err := run(args...)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return stdout, nil
}

func (cli *Cli) ListLabels(substr string) (bytes.Buffer, error) {
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

	args := []string{
		"graphql",
		"--paginate",
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
		"-F", fmt.Sprintf("label=%s", substr),
		"-f", fmt.Sprintf("query=%s", query),
	}

	stdout, _, err := run(args...)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return stdout, nil
}

func (cli *Cli) DeleteLabel(name string) error {
	args := []string{
		"-X", "DELETE",
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
		fmt.Sprintf("/repos/:owner/:repo/labels/%s", name),
	}

	_, _, err := run(args...)
	return err
}

func run(args ...string) (stdout, stderr bytes.Buffer, err error) {
	bin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("cannot find gh; is it installed? error: %w", err)
		return
	}

	// Always prepend arguments passed to every command.
	args = append([]string{"api", "-H", "accept:application/vnd.github.v3+json"}, args...)

	cmd := exec.Command(bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("gh returned error: %w, stderr: %s", err, stderr.String())
		return
	}

	return
}