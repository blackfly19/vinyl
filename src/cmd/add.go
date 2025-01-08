package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds file to the staging area",
	RunE:  addfunc,
}

func addfunc(cmd *cobra.Command, args []string) error {
	if args[0] == "." {
		files, err := os.ReadDir(".")
		if err != nil {
			return err
		}
		
		for _, file := range files{
			os.WriteFile(".svc/.staging",[]byte(file.Name()+"\n"), os.ModePerm)
		}
	}
	for _,file := range args{
		if _, err := os.Stat(file); err == nil{
			os.WriteFile(".svc/.staging",[]byte(file+"\n"),os.ModePerm)
		} else {
			fmt.Println("File not found.")
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
