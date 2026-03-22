package cmd

import (
	"github.com/blackfly19/godiff/diff"
	"github.com/blackfly19/vcs/src/constants"
	"os"
	"path/filepath"
	"reflect"

	"github.com/blackfly19/vcs/src/qwe"
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

		err = commit(message, checkpoint)
		if err != nil {
			return err
		}

		return nil
	},
}

func commit(message string, checkpoint bool) error {

	var commitID string
	stagingMap := qwe.NewMapHandler[qwe.FileMetaData](constants.FILE_STAGING)
	fileMap := qwe.NewMapHandler[qwe.FileMetaData](constants.FILE_FILEHASH)
	diffFileMap := qwe.NewMapHandler[string](constants.DIR_DIFF_MAP)
	tree := qwe.NewTree()
	stateTree := qwe.NewMerkleTree()

	err := stagingMap.ReadFromDisk()
	if err != nil {
		return err
	}

	err = fileMap.ReadFromDisk()
	if err != nil {
		return err
	}

	err = tree.ReadFromDisk()
	if err != nil {
		return err
	}

	stateTree.CreateSnapshot(".")

	commitID = qwe.GenerateCommitID(stateTree.Root.CompHash, message)

	diffFileMap.FilePath = diffFileMap.FilePath + commitID

	for path, fileMetaData := range stagingMap.FileMap {

		if reflect.ValueOf(fileMetaData).IsZero() {
			delete(fileMap.FileMap, path)
		} else {
			updatedFile, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = qwe.CreateDataObjects(updatedFile, fileMetaData.FileContentMD5Hash)
			if err != nil {
				return err
			}

			if !checkpoint && tree.Root != nil {
				var originalFile []byte
				var objectFilePath = constants.DIR_OBJECTS + fileMap.FileMap[path].FileContentMD5Hash

				objectFile, err := os.Open(objectFilePath)

				err = qwe.GobDecoder(objectFile, &originalFile)
				if err != nil {
					return err
				}

				diffFile := diff.Encode(updatedFile, originalFile, 8)
				diffFilePath := constants.DIR_DIFF + filepath.Base(objectFile.Name())

				err = os.WriteFile(diffFilePath, diffFile, 0644)
				if err != nil {
					return err
				}

				diffFileMap.FileMap[diffFilePath] = path
			}
		}

		fileMap.FileMap[path] = fileMetaData
	}

	tree.AddCommitToTree(commitID, checkpoint, message)
	err = tree.WriteToDisk()
	if err != nil {
		return err
	}

	err = fileMap.WriteToDisk()
	if err != nil {
		return err
	}

	err = os.Remove(stagingMap.FilePath)
	if err != nil {
		return err
	}

	err = stateTree.WriteToDisk(commitID)
	if err != nil {
		return err
	}

	err = diffFileMap.WriteToDisk()
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
