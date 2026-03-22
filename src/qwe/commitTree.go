package qwe

import (
	"github.com/blackfly19/vcs/src/constants"
	"os"
	"time"
)

type CommitTree struct {
	Root         *Commit
	HeadCommitID string
}

func NewTree() *CommitTree {
	return new(CommitTree)
}

func (tree *CommitTree) AddCommitToTree(commitID string, checkpoint bool, message string) {

	newNode := &Commit{CommitID: commitID, Checkpoint: checkpoint, Datetime: time.Now().String(), Message: message} // ParentNodeAddress: nil}

	if tree.Root == nil {
		tree.Root = newNode
		tree.HeadCommitID = commitID
	} else {
		head := AssignHead(tree.Root, tree.HeadCommitID)
		head.ChildNodeAddress = append(head.ChildNodeAddress, newNode)
		head = newNode
		tree.HeadCommitID = commitID
	}
}

func (tree *CommitTree) WriteToDisk() error {
	binary, err := os.OpenFile(constants.FILE_COMMITTREE, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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

func (tree *CommitTree) ReadFromDisk() error {
	binary, err := os.Open(constants.FILE_COMMITTREE)
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

func RebuildPointers(root *Commit, headCommitID string) *Commit {

	var node *Commit
	if root == nil {
		return nil
	}

	if root.CommitID == headCommitID {
		return root
	}

	for _, child := range root.ChildNodeAddress {
		if node == nil {
			node = RebuildPointers(child, headCommitID)
		}
		child.parentNodeAddress = root
	}

	if node != nil {
		return node
	}

	return nil
}

func AssignHead(root *Commit, headCommitID string) *Commit {
	if root == nil {
		return nil
	}

	if root.CommitID == headCommitID {
		return root
	}

	for _, child := range root.ChildNodeAddress {
		node := AssignHead(child, headCommitID)
		if node != nil {
			return node
		}
	}

	return nil
}
