package cmd

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os"

	"github.com/blackfly19/godiff/diff"
	"github.com/blackfly19/vcs/src/qwe"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Creates a new commit for the files",
	RunE:  commit,
}

func commit(cmd *cobra.Command, args []string) error {

	stagingFile, err := os.Open(".qwe/staging.gob")
	if err != nil {
		return err
	}
	defer func(stagingFile *os.File) {
		err := stagingFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stagingFile)

	stagingFileMap := make(map[string]qwe.FileMetaData)
	fileHashMap := make(map[string]bool)
	diffFileHashMap := make(map[string]bool)

	err = qwe.GobDecoder(stagingFile, &stagingFileMap)
	if err != nil {
		return err
	}

	tree := qwe.ReadTreeFromDisk(".qwe/committree.gob")

	filePathMetadataMap, err := qwe.ReadFileHashFromDisk(".qwe/filehashes.gob")
	if err != nil {
		return err
	}

	checkpoint, err := cmd.Flags().GetBool("checkpoint")
	if err != nil {
		return err
	}

	message, err := cmd.Flags().GetString("message")
	if err != nil {
		return err
	}

	for path, fileMetaData := range stagingFileMap {

		updatedFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Fix diff file logic
		if checkpoint && tree.Root != nil {
			var originalFile *bytes.Buffer
			var objectFilePath = ".qwe/objects/" + filePathMetadataMap[path].FileContentMD5Hash + ".gob"

			objectFile, err := os.Open(objectFilePath)

			err = qwe.GobDecoder(objectFile, originalFile)
			if err != nil {
				return err
			}

			diffFile := diff.Encode(originalFile.Bytes(), updatedFile, 8)
			diffFileName := fmt.Sprintf("%x", md5.Sum(diffFile))

			err = os.WriteFile(".qwe/diffobjects/"+diffFileName, diffFile, 0644)
			if err != nil {
				return err
			}

			diffFileHashMap[diffFileName] = true
		}

		err = qwe.CreateDataObjects(updatedFile, fileMetaData.FileContentMD5Hash)
		if err != nil {
			return err
		}

		filePathMetadataMap[path] = fileMetaData
		fileHashMap[fileMetaData.FileContentMD5Hash] = true
	}

	/*if tree.Root != nil {
		tree.Head.DiffFileMap = diffFileHashMap
	}*/

	tree.AddCommitToTree(fileHashMap, diffFileHashMap, checkpoint, message)
	qwe.WriteTreeToDisk(tree)

	err = qwe.WriteFileHashToDisk(".qwe/filehashes.gob", filePathMetadataMap)
	if err != nil {
		return err
	}

	err = os.Remove(stagingFile.Name())
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
