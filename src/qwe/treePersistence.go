package qwe

import (
	"log"
	"os"
)

func WriteTreeToDisk(tree *CommitTree) {
	binary, err := os.OpenFile(".qwe/committree.gob", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Unable to save the commit tree.")
	}
	defer func(binary *os.File) {
		err := binary.Close()
		if err != nil {
			log.Fatal("Unable to close the commit tree file.")
		}
	}(binary)

	tree.Head = nil

	err = GobEncoder(binary, &tree)
	if err != nil {
		//return error
	}

}

func ReadTreeFromDisk(fileName string) *CommitTree {
	binary, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to open the commit tree file.")
	}
	defer func(binary *os.File) {
		err := binary.Close()
		if err != nil {
			log.Fatal("Unable to close the commit tree file.")
		}
	}(binary)

	var tree CommitTree
	err = GobDecoder(binary, &tree)

	//addParentPointers(tree.Root, nil)
	if err != nil { // Return err
	}
	return &tree
}

/*func addParentPointers(node *Commit, parent *Commit) {
	if node == nil {
		return
	}

	node.ParentNodeAddress = parent
	for _, childNode := range node.ChildNodeAddress {
		addParentPointers(childNode, node)
	}
}*/
