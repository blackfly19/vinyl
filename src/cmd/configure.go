package cmd

/*
import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Initial configuration",
	RunE:  configureFunc,
}

func configureFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("configure called")
	return nil
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringP("global", "g", "false", "Setting global configuration")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	os.Stat(homeDir + "/")
}*/
