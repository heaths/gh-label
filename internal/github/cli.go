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
		"/repos/{owner}/{repo}/labels",
		"-X", "POST",
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
		"-F", fmt.Sprintf("name=%s", label.Name),
		"-f", fmt.Sprintf("color=%s", label.Color),
	}

	if label.Description != "" {
		args = append(args, "-F", fmt.Sprintf("description=%s", label.Description))
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
		fmt.Sprintf("/repos/{owner}/{repo}/labels/%s", name),
		"-X", "DELETE",
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
	}

	_, _, err := run(args...)
	return err
}

func (cli *Cli) UpdateLabel(label EditLabel) (bytes.Buffer, error) {
	args := []string{
		fmt.Sprintf("/repos/{owner}/{repo}/labels/%s", label.Name),
		"-X", "PATCH",
		"-F", fmt.Sprintf("owner=%s", cli.Owner),
		"-F", fmt.Sprintf("repo=%s", cli.Repo),
	}

	if label.Color != "" {
		args = append(args, "-f", fmt.Sprintf("color=%s", label.Color))
	}

	if label.Description != "" {
		args = append(args, "-F", fmt.Sprintf("description=%s", label.Description))
	}

	if label.NewName != "" {
		args = append(args, "-F", fmt.Sprintf("new_name=%s", label.NewName))
	}

	stdout, _, err := run(args...)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return stdout, nil
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
