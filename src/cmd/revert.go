package cmd

/*
import (
	"errors"
	"github.com/blackfly19/godiff/diff"
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/blackfly19/vcs/src/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Moves to a previous commit",
	RunE:  revert,
}

func revert(cmd *cobra.Command, args []string) error {

	var files []string
	currentStateTree := vinyl.NewMerkleTree()
	commitID := args[0]

	if commitID == "" {
		return errors.New("No commitID provided")
	}

	node, err := currentStateTree.ReadFromDisk(commitID)
	if err != nil {
		return err
	}

	for _, path := range files {
		path, _ = strings.CutPrefix(path, "/./")
		objectFilePath := currentStateTree.GetObjectFile(path)
		if objectFilePath != "" {
			var fileContent []byte
			file, err := os.Open(objectFilePath)
			if err != nil {
				return err
			}
			err = utils.GobDecoder(file, &fileContent)
			if err != nil {
				return err
			}
			err = os.WriteFile(path, fileContent, 0644)
			if err != nil {
				return err
			}
		} /*else {

			// Resolving deltas
			var originalFileContent []byte
			updatedFileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for iterCommit := head; iterCommit != targetCommit; iterCommit = iterCommit.GetParentNodeAddress() {
				diffFileMapPath := constants.DIR_DIFF_MAP + iterCommit.CommitID
				diffMap := vinyl.NewMapHandler[string](diffFileMapPath)
				if _, exists := diffMap.FileMap[path]; exists {
					diffFileContent, err := os.ReadFile(diffMap.FileMap[path])
					if err != nil {
						return err
					}
					originalFileContent = diff.Decode(updatedFileContent, diffFileContent)
					updatedFileContent = originalFileContent
				}

			}

			if len(originalFileContent) > 0 {
				err = os.WriteFile(path, originalFileContent, 0644)
				if err != nil {
					return err
				}
			} else {
				err = os.Remove(path)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(revertCmd)
	revertCmd.Flags().StringP("commitid", "c", "", "Specify which commit to pick up for restoring a file")
}*/
