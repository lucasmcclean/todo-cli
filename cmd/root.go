package cmd

import (
	"fmt"
	"os"

	"github.com/ljmcclean/todo-cli/menu"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	list        string
	create      bool
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a tool for managing tasks from the terminal",
	Long:  `Todo is a command line utility for generating and managing todo lists.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if list == "" {
			return fmt.Errorf("You must specifiy a list with the -l flag")
		}
		m, err := menu.New(list, create)
		if err != nil {
			return err
		}
		if interactive {
			m.RunInteractive()
		} else {
			m.PrintItems()
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&list, "list", "l", "", "specifiy list to display")
	rootCmd.Flags().BoolVarP(&create, "create", "c", false, "create list if it does not exist")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "make list interactive")
}
