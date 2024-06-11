package menu

import (
	"bufio"
	"fmt"
	"io"
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
	task    string
}

func (m *Menu) pushAction(op string, lineNum int, task string) {
	action := action{
		op:      op,
		lineNum: lineNum,
		task:    task,
	}
	m.history = append(m.history, action)
}

func (m *Menu) popAction() (action action) {
	action = m.history[len(m.history)-1]
	m.history = append(m.history, m.history[:len(m.history)-1]...)
	return action
}

func (m *Menu) MoveCursor(delta int) {
	m.pushAction("MOV", m.cursorPos, "")
	m.cursorPos += delta
	fileLen := GetFileLength(*m.file)
	for m.cursorPos < 0 {
		m.cursorPos += fileLen
	}
	for m.cursorPos >= fileLen {
		m.cursorPos -= fileLen
	}
}

func (m *Menu) PrintItems(isInteractive bool) {
	m.file.Seek(0, io.SeekStart)
	reader := bufio.NewReader(m.file)
	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if lineNum == m.cursorPos && isInteractive {
			fmt.Print(" > ")
		} else {
			fmt.Printf(" %d ", lineNum+1)
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
	if isInteractive {
		fmt.Println("Press 'h' for help")
	}
}
