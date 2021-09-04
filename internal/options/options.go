package options

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type keyStore interface {
	get(key string) string
}

type GlobalOptions struct {
	owner string
	repo  string

	// test
	keys keyStore
}

func New(cmd *cobra.Command) *GlobalOptions {
	opts := &GlobalOptions{}
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		repoOverride, _ := cmd.Flags().GetString("repo")
		return opts.parseRepoOverride(repoOverride)
	}

	cmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` format")

	return opts
}

func (opts *GlobalOptions) Repo() (owner, repo string) {
	return opts.owner, opts.repo
}

func (opts *GlobalOptions) parseRepoOverride(repoOverride string) error {
	if len(repoOverride) == 0 {
		if opts.keys == nil {
			opts.keys = &environment{}
		}
		repoOverride = opts.keys.get("GH_REPO")
	}

	if len(repoOverride) == 0 {
		opts.owner = ":owner"
		opts.repo = ":repo"
		return nil
	}

	parts := strings.Split(repoOverride, "/")

	if len(parts) != 2 {
		return fmt.Errorf(`expected the "OWNER/REPO" format, got %s`, repoOverride)
	}

	for _, part := range parts {
		if part == "" {
			return fmt.Errorf(`expected the "OWNER/REPO" format, got %s`, repoOverride)
		}
	}

	opts.owner = parts[0]
	opts.repo = parts[1]
	return nil
}

type environment struct{}

func (env *environment) get(key string) string {
	return os.Getenv(key)
}
