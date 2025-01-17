package cmd

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes new repository",
	RunE:  initfunc,
}

func initfunc(cmd *cobra.Command, args []string) error {

	err := os.Mkdir(".svc", os.ModePerm)
	if err != nil {

		return err
	}

	fmt.Println("New empty repository initialized.")
	initSnap, err := os.Create(".svc/v1.zip")
	if err != nil {
		return err
	}
	defer initSnap.Close()

	zipWriter := zip.NewWriter(initSnap)
	defer zipWriter.Close()

	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == ".svc" {
			return filepath.SkipDir
		}

		file, err := os.ReadFile(path)

		if !info.IsDir() {
			writer, err := zipWriter.Create(strings.TrimPrefix(path, "./"))
			if err != nil {
				return err
			}

			_, err = writer.Write(file)
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
