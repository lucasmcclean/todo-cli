package cmd

import (
	"fmt"

	"github.com/ljmcclean/todo-cli/menu"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Prints all existing todo lists",
	Long: `Accesses all files stored in the data directory for this todo
app (typically located in $HOME/.local/share/todo/) and prints them.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileNames, err := menu.GetDataFileNames()
		if err != nil {
			return nil
		}
		for _, fileName := range fileNames {
			fmt.Println(" - ", fileName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
