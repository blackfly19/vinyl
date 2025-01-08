package cmd

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes new repository",
	RunE:  initfunc,
}

func initfunc(cmd *cobra.Command, args []string) error {
	initSnap, err := os.Create("v1.zip")
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(initSnap)
	walkfn = filepath.Walk("")

	filepath.Walk("./",)

	err = os.Mkdir(".svc",os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println("New empty repository initialized.")
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
