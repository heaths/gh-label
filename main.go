package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cli/safeexec"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	owner string
	repo  string
}

func main() {
	opts := rootOptions{}

	rootCmd := cobra.Command{
		Use: "label",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			repoOverride, _ := cmd.Flags().GetString("repo")
			return opts.parseRepoOverride(repoOverride)
		},
	}

	rootCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` format")

	rootCmd.AddCommand(listCmd(&opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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

func (r *rootOptions) parseRepoOverride(repoOverride string) error {
	if len(repoOverride) == 0 {
		repoOverride = os.Getenv("GH_REPO")
	}

	if len(repoOverride) == 0 {
		return nil
	}

	parts := strings.SplitN(repoOverride, "/", 2)

	if len(parts) != 2 {
		return fmt.Errorf(`expected the "OWNER/REPO" format, got %s`, repoOverride)
	}

	for _, part := range parts {
		if part == "" {
			return fmt.Errorf(`expected the "OWNER/REPO" format, got %s`, repoOverride)
		}
	}

	r.owner = parts[0]
	r.repo = parts[1]
	return nil
}

func (r *rootOptions) repoOverride() (owner, repo string) {
	if owner = r.owner; owner == "" {
		owner = ":owner"
	}
	if repo = r.repo; repo == "" {
		repo = ":repo"
	}
	return
}
