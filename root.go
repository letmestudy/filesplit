package main

import "github.com/spf13/cobra"

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "file-splitter",
		Short: "A tool for splitting and merging files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(NewSplitCommand())
	rootCmd.AddCommand(NewMergeCommand())
	rootCmd.AddCommand(NewVerifyCommand())

	return rootCmd.Execute()
}
