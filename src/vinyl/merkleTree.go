package vinyl

import (
	"errors"
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/utils"
	"os"
	"path/filepath"
	"strings"
)

type MerkleTree struct {
	Root *ParentNode
}

func NewMerkleTree() *MerkleTree {
	return new(MerkleTree)
}

func (tree *MerkleTree) CreateSnapshot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tree.Root = &ParentNode{}
	err = BuildMerkleTreeNodes(tree.Root, cwd)
	if err != nil {
		return "", err
	}

	dirHash := tree.Root.CalculateObjectMD5()

	err = tree.Root.WriteToDisk(dirHash)
	if err != nil {
		return "", err
	}

	return dirHash, nil
}

func (tree *MerkleTree) GetObjectFile(path string) (string, error) {
	var err error
	var dataObjectFile string
	traversal := tree.Root

	dirs := strings.Split(path, "/")

	for _, dir := range dirs {
		var i int
		for i = 0; i < len(traversal.ChildNodes) && traversal.ChildNodes[i].Name != dir; i++ {
		}
		if i == len(traversal.ChildNodes) {
			return "", errors.New("Object not found")
		}
		if !traversal.ChildNodes[i].IsFile {
			traversal, err = tree.ReadFromDisk(traversal.ChildNodes[i].CompHash)
			if err != nil {
				return "", err
			}
		} else {
			dataObjectFile = traversal.ChildNodes[i].CompHash
		}
	}

	if _, err := os.Stat(constants.DIR_OBJECTS + dataObjectFile); err != nil {
		return "", errors.New("Object not found")
	} else {
		return constants.DIR_OBJECTS + dataObjectFile, nil
	}

}

func (tree *MerkleTree) ReadFromDisk(hash string) (*ParentNode, error) {
	var parentNode ParentNode
	gobfile, err := os.Open(constants.DIR_OBJECTS + hash)
	if err != nil {
		return nil, err
	}
	defer gobfile.Close()

	err = utils.GobDecoder(gobfile, &parentNode)
	if err != nil {
		return nil, err
	}

	return &parentNode, nil
}

func BuildMerkleTreeNodes(root *ParentNode, previousPath string) error {

	files, err := os.ReadDir(previousPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name()+"/" == constants.DIR_QWE {
			continue
		}
		if file.IsDir() {
			newNode := &ParentNode{}
			err = BuildMerkleTreeNodes(newNode, filepath.Join(previousPath, file.Name()))
			if err != nil {
				return err
			}
			dirHash := newNode.CalculateObjectMD5()
			err = newNode.WriteToDisk(dirHash)
			if err != nil {
				return err
			}
			dirChildNode := ChildNode{Name: file.Name(), IsFile: false, FileMode: file.Type().String(), CompHash: dirHash}
			root.ChildNodes = append(root.ChildNodes, dirChildNode)
		} else {
			newNode := ChildNode{Name: file.Name(), IsFile: true, FileMode: file.Type().String()}
			fileContent, err := os.ReadFile(filepath.Join(previousPath, file.Name()))
			if err != nil {
				return err
			}
			newNode.CompHash = utils.CalculateMD5(fileContent)
			err = CreateDataObjects(fileContent, newNode.CompHash)
			if err != nil {
				return err
			}
			root.ChildNodes = append(root.ChildNodes, newNode)
		}
	}

	return nil
}

func CreateDataObjects(fileContent []byte, fileName string) error {

	file, err := os.Create(constants.DIR_OBJECTS + fileName)
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

	err = utils.GobEncoder(file, fileContent)
	if err != nil {
		return err
	}

	return nil

}
