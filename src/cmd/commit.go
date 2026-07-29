package cmd

import (
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Creates a new commit for the files",
	RunE: func(cmd *cobra.Command, args []string) error {
		checkpoint, err := cmd.Flags().GetBool("checkpoint")
		if err != nil {
			return err
		}

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		err = commit(cmd, message, checkpoint)
		if err != nil {
			return err
		}

		return nil
	},
}

func commit(cmd *cobra.Command, message string, checkpoint bool) error {

	stateTree := vinyl.NewMerkleTree()
	commitTree, err := vinyl.LoadCommitTree()
	if err != nil {
		return err
	}

	rootHash, err := stateTree.CreateSnapshot()
	if err != nil {
		return err
	}

	err = commitTree.AddCommit(rootHash, checkpoint, message)
	if err != nil {
		return err
	}

	err = commitTree.WriteToDisk()
	if err != nil {
		return err
	}

	return nil

}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolP("checkpoint", "c", false, "To always keep the objects")
	commitCmd.Flags().StringP("message", "m", "No commit message", "Adds a message to the commit")
}
