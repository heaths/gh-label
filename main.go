package main

import (
	"os"

	"github.com/heaths/gh-label/internal/cmd/create"
	"github.com/heaths/gh-label/internal/cmd/list"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "label",
	}

	opts := options.New(&rootCmd)

	rootCmd.AddCommand(create.CreateCmd(opts))
	rootCmd.AddCommand(list.ListCmd(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
