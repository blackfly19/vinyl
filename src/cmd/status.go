package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/utils"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Lists all the modified files",
	RunE:  status,
}

func status(cmd *cobra.Command, args []string) error {

	var added, deleted, modified []string
	commitTree, err := vinyl.LoadCommitTree()
	if err != nil {
		return err
	}

	if commitTree == nil || commitTree.Head == "" {
		fmt.Println("Changes to be committed:")
		fmt.Println("Added files: ")

		err := filepath.WalkDir(".", func(path string, fde fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if fde.IsDir() && fde.Name() == ".vinyl" {
				return filepath.SkipDir
			}

			if !fde.IsDir() {
				fmt.Println(path)
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}

	targetCommit, err := vinyl.LoadCommit(commitTree.Head)
	if err != nil {
		return err
	}

	commitMerkleTree, err := vinyl.LoadMerkleTreeNode(targetCommit.MerkleRoot)
	if err != nil {
		return err
	}

	fmt.Println("Changes to be committed:")

	changeList := make(map[string]string)

	err = compareDirTree(targetCommit, commitMerkleTree, ".", changeList)
	if err != nil {
		return err
	}

	for key, value := range changeList {
		if value == "A" {
			added = append(added, key)
		} else if value == "D" {
			deleted = append(deleted, key)
		} else if value == "M" {
			modified = append(modified, key)
		} else if value == "DA" || value == "AD" {
			added = append(added, key)
			deleted = append(deleted, key)
		}
	}

	fmt.Println("Added files:", added)
	fmt.Println("Deleted files:", deleted)
	fmt.Println("Modified files:", modified)

	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func compareDirTree(targetCommit *vinyl.Commit, commitMerkleTree *vinyl.MerkleTree, nextWorkingDir string, changeList map[string]string) error {

	dirs, err := os.ReadDir(nextWorkingDir)
	if err != nil {
		return err
	}

	commitLen := len(commitMerkleTree.Root.ChildNodes)
	dirLen := len(dirs)

	currentNode := commitMerkleTree.Root

	i := 0
	j := 0
	for i < dirLen && j < commitLen {

		if dirs[i].Name()+"/" == constants.DIR_VINYL {
			i++
			continue
		}

		if dirs[i].Name() == currentNode.ChildNodes[j].Name {

			filePath := filepath.Join(nextWorkingDir, currentNode.ChildNodes[j].Name)

			if dirs[i].IsDir() && !currentNode.ChildNodes[j].IsFile {
				nextMerkleTreeNode, err := vinyl.LoadMerkleTreeNode(currentNode.ChildNodes[j].CompHash)
				if err != nil {
					return err
				}
				err = compareDirTree(targetCommit, nextMerkleTreeNode, nextWorkingDir+"/"+dirs[i].Name(), changeList)
				if err != nil {
					return err
				}
				// Condition to check file mode and if its different mark it
			} else if currentNode.ChildNodes[j].IsFile && !dirs[i].IsDir() {
				info, err := os.Stat(filePath)
				if err != nil {
					return err
				}
				commitTime, err := time.Parse(time.RFC3339, targetCommit.Datetime)
				if err != nil {
					return err
				}

				if info.ModTime().After(commitTime) {
					fileContent, err := os.ReadFile(filePath)
					if err != nil {
						return err
					}
					currentFileHash := utils.CalculateMD5(fileContent)
					if currentFileHash != currentNode.ChildNodes[j].CompHash {
						changeList[filePath] = "M"
					}
				}
				// Condition to check file mode and if it is different mark it
			} else if currentNode.ChildNodes[j].IsFile && dirs[i].IsDir() {
				changeList[filePath] = "DA"
			} else {
				changeList[filePath] = "AD"
			}
			i++
			j++

		} else if currentNode.ChildNodes[j].Name > dirs[i].Name() {
			// Mark that a new file/Dir has been added dirs[i].Name()
			changeList[filepath.Join(nextWorkingDir, dirs[i].Name())] = "A"
			i++
		} else {
			changeList[filepath.Join(nextWorkingDir, currentNode.ChildNodes[j].Name)] = "D"
			j++
		}

	}

	for ; i < dirLen; i++ {
		changeList[filepath.Join(nextWorkingDir, dirs[i].Name())] = "A"
	}

	for ; j < commitLen; j++ {
		changeList[filepath.Join(nextWorkingDir, currentNode.ChildNodes[j].Name)] = "D"
	}

	return nil
}
