package cmd

import (
	"errors"
	"github.com/blackfly19/vcs/src/utils"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restores a file to current commit state or a given commit",
	Args:  cobra.ExactArgs(1),
	RunE:  restore,
}

func restore(cmd *cobra.Command, args []string) error {
	var objectFileData []byte

	filePath := args[0]

	commitID, err := cmd.Flags().GetString("commit")
	if err != nil {
		return err
	}

	commitTree, err := vinyl.LoadCommitTree()
	if err != nil {
		return err
	}

	if commitTree == nil || commitTree.Head == "" {
		return errors.New("No commits exist")
	}

	if commitID == "" {
		commitID = commitTree.Head
	}

	targetCommit, err := vinyl.LoadCommit(commitID)
	if err != nil {
		return err
	}

	if targetCommit == nil {
		return errors.New("No commit found")
	}

	merkleTree, err := vinyl.LoadMerkleTreeNode(targetCommit.MerkleRoot)
	if err != nil {
		return err
	}

	objectFile, err := merkleTree.GetObjectFile(filePath)
	if err != nil {
		return err
	}

	file, err := os.Open(objectFile)
	if err != nil {
		return err
	}
	defer file.Close()

	err = utils.GobDecoder(file, &objectFileData)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	err = os.WriteFile(filePath, objectFileData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringP("commit", "c", "", "commit to restore")
}
