package cmd

import (
	"errors"
	"os"

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
		file, err := os.OpenFile(Directory+ListName, os.O_CREATE, 0600)
		if err != nil {
			return errors.New("Could not create file " + ListName + " at " + Directory)
		}
		defer file.Close()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
