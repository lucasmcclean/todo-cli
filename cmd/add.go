package cmd

import (
	// "errors"
	// "os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list",
	Long: `Adds a task to the currently active list unless a different
list is specified with the -l tag.`,
	/*
		 	RunE: func(cmd *cobra.Command, args []string) error {
				item := args[0]
				file, err := os.OpenFile(Directory+ListName, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					return errors.New("Could not open file " + ListName)
				}
				defer file.Close()

				_, err = file.WriteString(item + "\n")
				if err != nil {
					return errors.New("Failed to write item to file " + ListName)
				}
				return nil
			},
	*/
}

func init() {
	rootCmd.AddCommand(addCmd)
}
