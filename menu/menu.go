package menu

import (
	"bufio"
	"fmt"
	"os"
)

type Menu struct {
	file      *os.File
	cursorPos int
	history   []action
}

func New(fileName string, create bool) (menu *Menu, err error) {
	file, err := OpenDataFile(fileName, create)
	if err != nil {
		return nil, err
	}
	menu = &Menu{
		file:      file,
		cursorPos: 0,
		history:   make([]action, 20),
	}
	return menu, nil
}

type action struct {
	op      string
	lineNum int
}

func (m *Menu) PrintItems() {
	reader := bufio.NewReader(m.file)
	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if len(line) <= 1 {
			continue
		}
		if lineNum == m.cursorPos {
			fmt.Print(" > ")
		} else {
			fmt.Print("   ")
		}
		if len(line) >= 5 && line[:5] == "[!-!]" {
			fmt.Print("[X] ")
			line = line[5:]
		} else {
			fmt.Print("[ ] ")
		}
		fmt.Print(line)
		lineNum++
	}
}
