package cmd

import (
	"os"

	utils "github.com/ljmcclean/todo-cli/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	Verbose   bool
	Directory string
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a tool for managing tasks from the terminal",
	Long: `Todo is a command line utility for generating and managing
todo lists.
A call to the root command lists the tasks from your active
todo list unless a seperate path is specified with -l.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return drawTodo()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	dir := utils.MustGetDataDir()

	rootCmd.PersistentFlags().StringVarP(&Directory, "directory", "d", dir, "specify directory to store lists")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display more verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func drawTodo() error {
	restore, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), restore)

	in := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(in)
		if err != nil {
			return err
		}
		if in[0] == 113 {
			break
		}
	}
	return nil
}

func renderTodo() {

}
