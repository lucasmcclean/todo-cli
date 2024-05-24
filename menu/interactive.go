package menu

import (
	"os"

	"golang.org/x/term"
)

const (
	Quit   byte = 113
	Down   byte = 107
	UpA    byte = 183
	Up     byte = 106
	DownA  byte = 184
	Del    byte = 120
	Mark   byte = 13
	After  byte = 97
	Insert byte = 105
	Make   byte = 109
	Undo   byte = 117
)

func (m *Menu) RunInteractive() (err error) {
Interactive:
	for {
		inputCode, err := getRawInput()
		if err != nil {
			return err
		}
		switch inputCode {
		case Quit:
			break Interactive
		case Down, DownA:
		case Up, UpA:
		case Mark:
		case Del:
		case After:
		case Insert:
		case Undo:
		case Make:
		}
	}
	return nil
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
