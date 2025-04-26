package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/qwe"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes new repository",
	RunE:  initFunc,
}

func initFunc(cmd *cobra.Command, args []string) error {

	// Create hidden directory for qwe
	err := os.Mkdir(".qwe", os.ModePerm)
	if err != nil {
		return err
	}

	//Create object directory to store objects
	err = os.Mkdir(".qwe/objects", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".qwe/diffobjects", os.ModePerm)
	if err != nil {
		return err
	}

	//Initializing file hashes
	err = qwe.WriteFileHashToDisk(".qwe/filehashes.gob", make(map[string]qwe.FileMetaData))
	if err != nil {
		return err
	}

	//Initializing commit tree
	qwe.InitializeTree()

	fmt.Println("New empty repository initialized.")

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
