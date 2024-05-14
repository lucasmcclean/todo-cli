package cmd

import (
	"fmt"
	"os"

	utils "github.com/ljmcclean/todo-cli/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	Verbose   bool
	Directory string

	offset int
	remove bool
	menu   *utils.Menu
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a tool for managing tasks from the terminal",
	Long: `Todo is a command line utility for generating and managing
todo lists.
A call to the root command lists the tasks from your active
todo list unless a seperate path is specified with -l.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		menu = utils.NewMenu(Directory + args[0])
		defer menu.Close()
		return runInteractiveTodo()
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

	rootCmd.Flags().IntVarP(&offset, "offset", "o", 1, "number of lines to keep free above list")
	viper.BindPFlag("offset", rootCmd.Flags().Lookup("offset"))

	rootCmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove items that are checked")

	rootCmd.PersistentFlags().StringVarP(&Directory, "directory", "d", dir, "specify directory to store lists")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display more verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func runInteractiveTodo() error {
DrawLoop:
	for {
		menu.Draw(offset)
		inputCode, err := handleInput()
		if err != nil {
			panic(err)
		}
		switch inputCode {
		case utils.Q:
			break DrawLoop
		case utils.J, utils.Da:
			menu.MoveCursor(-1)
		case utils.K, utils.Ua:
			menu.MoveCursor(1)
		case utils.En:
			menu.CompleteItem(remove)
		case utils.X:
			menu.CompleteItem(true)
		default:
			fmt.Println(inputCode)
			return nil
		}
	}
	return nil
}

func handleInput() (inputCode byte, err error) {
	restore, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 1, err
	}
	defer term.Restore(int(os.Stdin.Fd()), restore)

	input := make([]byte, 1)
	_, err = os.Stdin.Read(input)
	if err != nil {
		return 1, err
	}
	return input[0], nil
}
