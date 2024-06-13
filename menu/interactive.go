package menu

import (
	"bufio"
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
	Move   byte = 109
	Undo   byte = 117
	Help   byte = 104
)

func (m *Menu) RunInteractive() (err error) {
	helpMenu := drawHelp()
	isMoving := false
	oldPos := 0
Interactive:
	for {
		taskList := m.DrawMenu(true)
		clearScreen()
		fmt.Print(taskList)
		if isMoving {
			fmt.Println("Moving item...\nPress 'a' to place after selected item or 'i' to place before\nPress'q' to cancel")
		}
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
			m.MarkItem(m.cursorPos)
		case Del:
			m.DeleteItem(m.cursorPos)
		case After:
			if isMoving {
				m.MoveItem(oldPos, m.cursorPos+1)
				isMoving = false
			} else {
				usrIn, err := promptItem()
				if err != nil {
					return err
				}
				m.CreateItem(m.cursorPos+1, usrIn)
			}
		case Insert:
			if isMoving {
				m.MoveItem(oldPos, m.cursorPos)
				isMoving = false
			} else {
				usrIn, err := promptItem()
				if err != nil {
					return err
				}
				m.CreateItem(m.cursorPos, usrIn)
			}
		case Undo:
		case Move:
			isMoving = true
			oldPos = m.cursorPos
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
		"Move the current task     '%s'\n"+
		"Undo the last action      '%s'\n"+
		"Open help menu            '%s'\n"+
		"-----------------------------------------\n"+
		"Inserting before ('i') and after ('a') creates a new item"+
		"Hit 'q' to exit this help menu\n",
		string(Quit), string(Down), string(Up),
		string(Del), string(After), string(Insert),
		string(Move), string(Undo), string(Help))
	return output
}

func promptItem() (userInput string, err error) {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print("New item: ")
	userInput, err = inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return userInput, nil
}
