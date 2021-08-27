package main

import (
	"os"

	"github.com/heaths/gh-label/internal/cmd/list"
	"github.com/heaths/gh-label/internal/options"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "label",
	}

	opts := options.RootOptions{}
	opts.Init(&rootCmd)

	rootCmd.AddCommand(list.ListCmd(&opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
