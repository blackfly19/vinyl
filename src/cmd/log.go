package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/vinyl"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays commit ids and messages",
	RunE:  logFunc,
}

func logFunc(cmd *cobra.Command, args []string) error {

	tree, err := vinyl.LoadCommitTree()
	if err != nil {
		return err
	}

	traversal := tree.Head

	for traversal != "" {
		treeCommit, err := tree.ReadFromDisk(traversal)
		if err != nil {
			return err
		}
		fmt.Println("Commit: ", treeCommit.CommitID)
		fmt.Println("Date: ", treeCommit.Datetime)
		fmt.Println("Message: ", treeCommit.Message)
		fmt.Println()
		traversal = treeCommit.ParentNodeCommitID
	}

	return nil
}

func init() {
	rootCmd.AddCommand(logCmd)
}
