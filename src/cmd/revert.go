package cmd

import (
	"errors"
	//"github.com/blackfly19/godiff/diff"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
)

var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Moves to a previous commit",
	RunE:  revert,
}

func revert(cmd *cobra.Command, args []string) error {

	commitID := args[0]

	if commitID == "" {
		return errors.New("No commitID provided")
	}

	//Load Commit

	node, err := vinyl.LoadCommit(commitID)
	if err != nil {
		return err
	}

	if node == nil {
		return errors.New("No commit found")
	}

	_, err = vinyl.LoadMerkleTreeNode(node.MerkleRoot)
	if err != nil {
		return err
	}

	//merkleTree.

	return nil
}

func init() {
	rootCmd.AddCommand(revertCmd)
	revertCmd.Flags().StringP("commitid", "c", "", "Specify which commit to pick up for restoring a file")
}
