package cmd

import (
	"fmt"
	"github.com/blackfly19/vcs/src/constants"
	"github.com/blackfly19/vcs/src/vinyl"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes new store",
	RunE:  initFunc,
}

func initFunc(cmd *cobra.Command, args []string) error {

	// Create hidden directory for vinyl
	err := os.Mkdir(constants.DIR_VINYL, os.ModePerm)
	if err != nil {
		return err
	}

	//Create object directory to store objects
	err = os.Mkdir(constants.DIR_OBJECTS, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(constants.DIR_DIFF, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(constants.DIR_COMMITS, os.ModePerm)
	if err != nil {
		return err
	}

	commitTree := vinyl.NewCommitTree()
	err = commitTree.WriteToDisk()
	if err != nil {
		return err
	}

	//Initializing file hashes
	err = vinyl.NewMapHandler[vinyl.FileMetaData](constants.FILE_FILEHASH).WriteToDisk("")
	if err != nil {
		return err
	}

	fmt.Println("New empty store initialized.")

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
