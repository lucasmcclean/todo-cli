package cmd

import (
	"fmt"

	"github.com/ljmcclean/todo-cli/menu"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Deletes the specified todo list",
	Long: `Deletes the todo list with the specified name.
! This is not recoverable !`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		for input != "Y" && input != "n" {
			fmt.Println("Are you sure you want to permanently delete '" + args[0] + "'? (Y/n)")
			fmt.Scan(&input)
		}
		if input == "n" {
			return nil
		}
		err := menu.RemoveDataFile(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
