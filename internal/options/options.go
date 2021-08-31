package options

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type config interface {
	get(key string) string
}

type RootOptions struct {
	owner string
	repo  string
	env   config
}

func New(cmd *cobra.Command) *RootOptions {
	r := &RootOptions{}
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		repoOverride, _ := cmd.Flags().GetString("repo")
		return r.parseRepoOverride(repoOverride)
	}

	cmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` format")

	return r
}

func (r *RootOptions) ConfigureCommand(cmd *cobra.Command) {
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Parent().PersistentPreRunE != nil {
			return cmd.Parent().PersistentPreRunE(cmd, args)
		}
		return nil
	}
}

func (r *RootOptions) Repo() (owner, repo string) {
	return r.owner, r.repo
}

func (r *RootOptions) parseRepoOverride(repoOverride string) error {
	if len(repoOverride) == 0 {
		if r.env == nil {
			r.env = &envConfig{}
		}
		repoOverride = r.env.get("GH_REPO")
	}

	if len(repoOverride) == 0 {
		r.owner = ":owner"
		r.repo = ":repo"
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

	r.owner = parts[0]
	r.repo = parts[1]
	return nil
}

type envConfig struct{}

func (c *envConfig) get(key string) string {
	return os.Getenv(key)
}
