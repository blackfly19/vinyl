package cmd

/*
import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a vcs server",
	RunE:  serve,
}

func serve(cmd *cobra.Command, args []string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Qwe!")
	})

	port, err := cmd.Flags().GetString("port")
	if err != nil {
		return err
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}

	fmt.Println("Started server on http://127.0.0.1:" + port)

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Port number to serve on")
}*/
