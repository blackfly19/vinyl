package qwe

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Commit struct {
	CommitID      string
	CommitFileMap map[string]bool
	DiffFileMap   map[string]bool
	Datetime      string
	Checkpoint    bool
	Message       string
	//ParentNodeAddress *Commit `gob:"-"`
	ChildNodeAddress []*Commit
}

type CommitTree struct {
	Root *Commit
	Head *Commit
}

func GenerateCommitID() string {
	newHash := md5.Sum([]byte(time.Now().String()))
	return hex.EncodeToString(newHash[:])
}

func InitializeTree() {
	tree := new(CommitTree)
	WriteTreeToDisk(tree)
}

func (tree *CommitTree) AddCommitToTree(commitFileMap map[string]bool, diffFileMap map[string]bool, checkpoint bool, message string) {

	commitID := GenerateCommitID()
	fmt.Println(commitID)

	newNode := &Commit{CommitID: commitID, CommitFileMap: commitFileMap, DiffFileMap: diffFileMap, Checkpoint: checkpoint, Datetime: time.Now().String(), Message: message} // ParentNodeAddress: nil}

	if tree.Root == nil {
		tree.Root = newNode
		tree.Head = newNode
	} else {
		tree.Head = findHead(tree.Root, tree.Head.CommitID)
		tree.Head.ChildNodeAddress = append(tree.Head.ChildNodeAddress, newNode)
		tree.Head = newNode
	}
}

func findHead(root *Commit, head string) *Commit {

	if root == nil {
		return nil
	}

	if root.CommitID == head {
		return root
	}

	for _, child := range root.ChildNodeAddress {
		node := findHead(child, head)
		if node != nil {
			return node
		}
	}

	return nil
}
