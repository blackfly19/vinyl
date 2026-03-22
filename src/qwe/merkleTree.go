package qwe

import (
	"fmt"
	"github.com/blackfly19/vcs/src/constants"
	"os"
	"path/filepath"
	"strings"
)

type MerkleTree struct {
	Root *MerkleTreeNode
}

func NewMerkleTree() *MerkleTree {
	return new(MerkleTree)
}

func (tree *MerkleTree) CreateSnapshot(dirName string) {
	tree.Root = &MerkleTreeNode{Name: dirName, IsFile: false}
	err := BuildMerkleTree(tree.Root, dirName)
	if err != nil {
		panic(err)
	}
}

func (tree *MerkleTree) WriteToDisk(fileName string) error {
	binary, err := os.OpenFile(constants.DIR_STATE_TREE+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(binary *os.File) {
		err = binary.Close()
	}(binary)
	if err != nil {
		return err
	}

	err = GobEncoder(binary, tree)
	if err != nil {
		return err
	}

	return nil
}

func (tree *MerkleTree) ReadFromDisk(fileName string) error {
	binary, err := os.Open(constants.DIR_STATE_TREE + fileName)
	if err != nil {
		return err
	}
	defer func(binary *os.File) {
		err = binary.Close()
	}(binary)
	if err != nil {
		return err
	}

	err = GobDecoder(binary, tree)

	if err != nil {
		return err
	}

	return nil
}

func (tree *MerkleTree) GetObjectFile(path string) string {
	traversal := tree.Root

	dirs := strings.Split(path, "/")

	for _, dir := range dirs {
		var i int
		for i = 0; i < len(traversal.ChildNodes) && traversal.ChildNodes[i].Name != dir; i++ {
		}
		traversal = traversal.ChildNodes[i]
	}

	if _, err := os.Stat(constants.DIR_OBJECTS + traversal.DataObject); err != nil {
		return ""
	} else {
		return constants.DIR_OBJECTS + traversal.DataObject
	}

}

func BuildMerkleTree(root *MerkleTreeNode, previousPath string) error {
	if !root.IsFile {
		if root.Name+"/" == constants.DIR_QWE {
			return nil
		}
		files, err := os.ReadDir(previousPath)
		if err != nil {
			return err
		}

		for _, file := range files {
			newNode := &MerkleTreeNode{Name: file.Name()}
			fmt.Println(newNode.Name)
			if file.IsDir() {
				newNode.IsFile = false
				err = BuildMerkleTree(newNode, filepath.Join(previousPath, newNode.Name))
				if err != nil {
					return err
				}
			} else {
				newNode.IsFile = true
				err = newNode.CalculateObjectMD5(filepath.Join(previousPath, newNode.Name))
				if err != nil {
					return err
				}
			}
			root.ChildNodes = append(root.ChildNodes, newNode)
		}

		err = root.CalculateObjectMD5(previousPath)
		if err != nil {
			return err
		}

	}
	return nil
}

func GetFilePaths(root *MerkleTreeNode, files *[]string, path string) error {
	if root.IsFile {
		*files = append(*files, path+"/"+root.Name)
		return nil
	}

	for _, node := range root.ChildNodes {
		err := GetFilePaths(node, files, path+"/"+root.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
