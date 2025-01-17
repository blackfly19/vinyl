package cmd

import (
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Creates a new commit for the files",
	//RunE:  commitfunc,
}

//func commitfunc(cmd *cobra.Command, args []string) error { }

func init() {
	rootCmd.AddCommand(commitCmd)
}
