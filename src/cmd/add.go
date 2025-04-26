package cmd

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blackfly19/vcs/src/qwe"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds file to the staging area",
	RunE:  add,
}

func add(cmd *cobra.Command, args []string) error {

	projectFileHashes, err := qwe.ReadFileHashFromDisk(".qwe/filehashes.gob")
	if err != nil {
		return err
	}

	stagingFileMap := make(map[string]qwe.FileMetaData)

	if len(args) > 0 {
		if args[0] == "." {

			fmt.Println("Added the following files: ")

			err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() && info.Name() == ".qwe" {
					return filepath.SkipDir
				}

				if !info.IsDir() {
					path = strings.TrimPrefix(path, "./")
					err := FileChange(projectFileHashes, stagingFileMap, path, info)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			for _, filePath := range args {
				if _, exists := projectFileHashes[filePath]; exists {
					info, err := os.Stat(filePath)
					if err != nil {
						return err
					}

					err = FileChange(projectFileHashes, stagingFileMap, filePath, info)
					if err != nil {
						return err
					}
				} else {
					fmt.Println("File not found.")
				}
			}
		}
	} else {
		log.Fatal("No argument given")
	}

	if len(stagingFileMap) > 0 {
		stagingFile, err := os.Create(".qwe/staging.gob")
		if err != nil {
			return errors.New("repository not initialized or not in root folder")
		}
		defer func(stagingFile *os.File) error {
			err := stagingFile.Close()
			if err != nil {
				return err
			}
			return nil
		}(stagingFile)

		err = qwe.GobEncoder(stagingFile, stagingFileMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileChange(projectFileHashes map[string]qwe.FileMetaData, stagingFileMap map[string]qwe.FileMetaData, path string, info fs.FileInfo) error {
	if _, exists := projectFileHashes[path]; exists {
		if info.ModTime().After(projectFileHashes[path].LastModifiedTime) {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			updatedHash := fmt.Sprintf("%x", md5.Sum(file))

			if updatedHash != projectFileHashes[path].FileContentMD5Hash {
				fmt.Println(path)
				stagingFileMap[path] = qwe.FileMetaData{FileContentMD5Hash: updatedHash, LastModifiedTime: info.ModTime()}
			}
		}
	} else {
		stagingFileMap[path] = qwe.FileMetaData{FileContentMD5Hash: fmt.Sprintf("%x", md5.Sum([]byte(path))), LastModifiedTime: info.ModTime()}
		fmt.Println(path)
	}
	return nil
}

func checkIfFileDeleted(projectFileHashes map[string]qwe.FileMetaData) {

	for file, _ := range projectFileHashes {
		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			delete(projectFileHashes, file)
		}
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
}
