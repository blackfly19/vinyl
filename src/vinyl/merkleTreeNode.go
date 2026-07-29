package vinyl

import (
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/utils"
	"os"
)

type ChildNode struct {
	Name     string
	IsFile   bool
	FileMode string
	CompHash string
}

type ParentNode struct {
	ChildNodes []ChildNode
}

func (node *ParentNode) CalculateObjectMD5() string {

	dirDesc := ""
	for _, child := range node.ChildNodes {
		dirDesc = dirDesc + child.Name + " " + child.FileMode + " " + child.CompHash + " "
	}
	return utils.CalculateMD5(dirDesc)
}

func (node *ParentNode) WriteToDisk(fileName string) error {
	file, err := os.OpenFile(constants.DIR_OBJECTS+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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

	err = utils.GobEncoder(file, node)
	if err != nil {
		return err
	}

	return nil
}
