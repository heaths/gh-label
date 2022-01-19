package main

import (
	"os"

	"github.com/heaths/gh-label/internal/cmd/create"
	"github.com/heaths/gh-label/internal/cmd/delete"
	"github.com/heaths/gh-label/internal/cmd/edit"
	"github.com/heaths/gh-label/internal/cmd/export"
	"github.com/heaths/gh-label/internal/cmd/list"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "label",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	opts := options.New(&rootCmd)

	rootCmd.AddCommand(create.CreateCmd(opts))
	rootCmd.AddCommand(edit.EditCmd(opts))
	rootCmd.AddCommand(export.ExportCmd(opts))
	rootCmd.AddCommand(delete.DeleteCmd(opts))
	rootCmd.AddCommand(list.ListCmd(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
