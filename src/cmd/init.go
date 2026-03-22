package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/constants"
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
	err := os.Mkdir(constants.DIR_QWE, os.ModePerm)
	if err != nil {
		return err
	}

	//Create object directory to store objects
	err = os.Mkdir(constants.DIR_OBJECTS, os.ModePerm)
	if err != nil {
		return err
	}

	//
	err = os.Mkdir(constants.DIR_DIFF, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(constants.DIR_STATE_TREE, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(constants.DIR_DIFF_MAP, os.ModePerm)
	if err != nil {
		return err
	}

	//Initializing file hashes
	err = qwe.NewMapHandler[qwe.FileMetaData](constants.FILE_FILEHASH).WriteToDisk()
	if err != nil {
		return err
	}

	//Initializing commit tree
	err = qwe.NewTree().WriteToDisk()
	if err != nil {
		return err
	}

	fmt.Println("New empty repository initialized.")

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
