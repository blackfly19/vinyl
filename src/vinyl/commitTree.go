package vinyl

import (
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/utils"
	"os"
	"time"
)

type CommitTree struct {
	Head string
}

func NewCommitTree() *CommitTree {
	return &CommitTree{Head: ""}
}

func LoadCommitTree() (*CommitTree, error) {
	var commitTree CommitTree
	gobfile, err := os.Open(constants.FILE_HEAD)
	if err != nil {
		return nil, err
	}
	defer gobfile.Close()

	err = utils.GobDecoder(gobfile, &commitTree)
	if err != nil {
		return nil, err
	}

	return &commitTree, nil
}

func (ct *CommitTree) AddCommit(merkleRoot string, checkpoint bool, message string) error {
	currentDateTime := time.Now().Format(time.RFC3339)
	commitID := utils.GenerateCommitID(merkleRoot, message, currentDateTime)
	newCommit := Commit{CommitID: commitID, MerkleRoot: merkleRoot, Checkpoint: checkpoint, Message: message, ParentNodeCommitID: ct.Head, Datetime: currentDateTime}

	err := newCommit.WriteToDisk(commitID)
	if err != nil {
		return err
	}

	ct.Head = commitID

	return nil
}

func (ct *CommitTree) ReadFromDisk(fileName string) (*Commit, error) {
	var commit *Commit
	gobfile, err := os.Open(constants.DIR_COMMITS + fileName)
	if err != nil {
		return nil, err
	}
	defer gobfile.Close()

	err = utils.GobDecoder(gobfile, &commit)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (ct *CommitTree) WriteToDisk() error {
	file, err := os.OpenFile(constants.FILE_HEAD, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	err = utils.GobEncoder(file, ct)
	if err != nil {
		return err
	}

	return nil
}
