package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/qwe"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds file to the staging area",
	RunE:  add,
}

func add(cmd *cobra.Command, args []string) error {

	fileMap := qwe.NewMapHandler[qwe.FileMetaData](constants.FILE_FILEHASH)
	stagingMap := qwe.NewMapHandler[qwe.FileMetaData](constants.FILE_STAGING)
	err := fileMap.ReadFromDisk()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		if args[0] == "." {

			fmt.Println("Added the following files: ")

			err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() && info.Name() == ".qwe" {
					return filepath.SkipDir
				}

				if !info.IsDir() {
					path = strings.TrimPrefix(path, "./")
					metadata, err := qwe.IsModified(fileMap.FileMap, path)
					if err != nil {
						return err
					}

					if !reflect.ValueOf(metadata).IsZero() {
						fmt.Println(path)
						stagingMap.FileMap[path] = metadata
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			for _, filePath := range args {
				if _, exists := fileMap.FileMap[filePath]; exists {

					metadata, err := qwe.IsModified(fileMap.FileMap, filePath)
					if err != nil {
						return err
					}

					if !reflect.ValueOf(metadata).IsZero() {
						stagingMap.FileMap[filePath] = metadata
					}
				}
			}
		}
	} else {
		log.Fatal("No argument given")
	}

	deletedFiles := qwe.DeletedFiles(fileMap.FileMap)

	if len(deletedFiles) > 0 {
		fmt.Println("Deleted files: ")
		for _, file := range deletedFiles {
			delete(fileMap.FileMap, file)
			stagingMap.FileMap[file] = qwe.FileMetaData{}
		}
	}

	if len(stagingMap.FileMap) > 0 {
		err = stagingMap.WriteToDisk()
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
