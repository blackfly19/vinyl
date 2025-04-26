package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/qwe"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays commit ids and messages",
	RunE:  logFunc,
}

func logFunc(cmd *cobra.Command, args []string) error {
	tree := qwe.ReadTreeFromDisk(".qwe/committree.gob")

	treeTraversal(tree.Root)
	return nil
}

func treeTraversal(node *qwe.Commit) {
	if node == nil {
		return
	}

	fmt.Printf("%s %s", node.CommitID, node.Message)
	for _, childNode := range node.ChildNodeAddress {
		treeTraversal(childNode)
	}
}

func init() {
	rootCmd.AddCommand(logCmd)
}
