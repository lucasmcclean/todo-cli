package menu

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	Quit   byte = 113
	QuitC  byte = 3
	Down   byte = 106
	UpA    byte = 183
	Up     byte = 107
	DownA  byte = 184
	Del    byte = 120
	Mark   byte = 13
	After  byte = 97
	Insert byte = 105
	Make   byte = 109
	Undo   byte = 117
	Help   byte = 104
)

func (m *Menu) RunInteractive() (err error) {
	helpMenu := drawHelp()
Interactive:
	for {
		taskList := m.DrawMenu(true)
		clearScreen()
		fmt.Print(taskList)
		inputCode, err := getRawInput()
		if err != nil {
			return err
		}
		switch inputCode {
		case Quit, QuitC:
			break Interactive
		case Down, DownA:
			m.MoveCursor(1)
		case Up, UpA:
			m.MoveCursor(-1)
		case Mark:
			m.MarkItem()
		case Del:
		case After:
		case Insert:
		case Undo:
		case Make:
		case Help:
			clearScreen()
			fmt.Print(helpMenu)
			inputCode, err = getRawInput()
			if err != nil {
				return err
			}
			for inputCode != Quit && inputCode != QuitC {
				inputCode, err = getRawInput()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func clearScreen() {
	fmt.Print("\033[H")
	fmt.Print("\033[2J")
}

func getRawInput() (inputCode byte, err error) {
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

func drawHelp() (output string) {
	output = fmt.Sprintf(""+
		"Quit this session         '%s' or 'Ctrl-c'\n"+
		"Move cursor down          'Up' or '%s'\n"+
		"Move cursor up            'Down'  or '%s'\n"+
		"Delete current task       '%s'\n"+
		"Mark task complete        'Enter'\n"+
		"Insert task after cursor  '%s'\n"+
		"Insert task before cursor '%s'\n"+
		"Create a new task         '%s'\n"+
		"Undo the last action      '%s'\n"+
		"Open help menu            '%s'\n"+
		"-----------------------------------------\n"+
		"Hit 'q' to exit this help menu\n",
		string(Quit), string(Down), string(Up),
		string(Del), string(After), string(Insert),
		string(Make), string(Undo), string(Help))
	return output
}
