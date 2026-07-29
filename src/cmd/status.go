package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Lists all the modified files",
	RunE:  status,
}

func status(cmd *cobra.Command, args []string) error {

	fileMap := vinyl.NewMapHandler[vinyl.FileMetaData](constants.FILE_FILEHASH)

	err := fileMap.ReadFromDisk("")
	if err != nil {
		return err
	}

	fmt.Println("Untracked files:")

	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".vinyl" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			path = strings.TrimPrefix(path, "./")
			metadata, err := vinyl.IsModified(fileMap.FileMap, path)
			if err != nil {
				return err
			}

			if !reflect.ValueOf(metadata).IsZero() {
				fmt.Println(path)
			}
		}
		return nil
	})

	deletedFiles := vinyl.DeletedFiles(fileMap.FileMap)

	if len(deletedFiles) > 0 {
		fmt.Println("Deleted files: ")
		for _, file := range deletedFiles {
			fmt.Println(file)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
