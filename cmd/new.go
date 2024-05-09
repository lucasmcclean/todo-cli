package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new todo list",
	Long:  `Generates a new list that you can add tasks to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ListName == "_current" {
			return errors.New("List argument required (-l, --list), \"_current\" is not a valid list name")
		}
		fmt.Printf("new %s at %s\n", ListName, Directory)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
