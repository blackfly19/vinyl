package qwe

import (
	"io"
	"os"
)

type MerkleTreeNode struct {
	Name       string
	IsFile     bool
	FileMode   string
	CompHash   string
	DataObject string
	ChildNodes []*MerkleTreeNode
}

func (node *MerkleTreeNode) CalculateObjectMD5(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) error {
		err := file.Close()
		if err != nil {
			return err
		}
		return nil
	}(file)
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if node.IsFile {
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		node.FileMode = fileInfo.Mode().String()
		node.DataObject = CalculateMD5(fileContent)
		node.CompHash = CalculateMD5(node.Name + " " + node.FileMode + " " + node.DataObject)
	} else {
		dirDesc := ""
		for _, childNode := range node.ChildNodes {
			dirDesc = dirDesc + childNode.Name + " " + childNode.FileMode + " " + childNode.CompHash
		}
		node.FileMode = fileInfo.Mode().String()
		node.CompHash = CalculateMD5(dirDesc)
	}

	return nil
}
