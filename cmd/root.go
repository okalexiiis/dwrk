package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "dwrk",
	Short: "Project Management Tool",
}

func Execute() {
	RootCmd.Execute()
}
