package vinyl

import (
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/utils"
	"os"
)

type Commit struct {
	CommitID           string
	MerkleRoot         string
	Datetime           string
	Checkpoint         bool
	Message            string
	ParentNodeCommitID string
}

func (c *Commit) WriteToDisk(fileName string) error {
	file, err := os.OpenFile(constants.DIR_COMMITS+fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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

	err = utils.GobEncoder(file, c)
	if err != nil {
		return err
	}

	return nil
}
