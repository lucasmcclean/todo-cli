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
	rw        bufio.ReadWriter
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
		rw:        *bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file)),
	}
	return menu, nil
}

type action struct {
	op      string
	lineNum int
	task    string
}

func (m *Menu) newAction(op string, lineNum int, task string) {
	action := action{
		op:      op,
		lineNum: lineNum,
		task:    task,
	}
	m.history = append(m.history, action)
}

func (m *Menu) undoAction() (action action) {
	action = m.history[len(m.history)-1]
	m.history = append(m.history, m.history[:len(m.history)-1]...)
	return action
}

func (m *Menu) MoveCursor(delta int) {
	m.newAction("MOV", m.cursorPos, "")
	m.cursorPos += delta
	fileLen := GetFileLength(*m.file)
	for m.cursorPos < 0 {
		m.cursorPos += fileLen
	}
	for m.cursorPos >= fileLen {
		m.cursorPos -= fileLen
	}
}

func (m *Menu) DeleteItem() {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	var deletedItem string
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == m.cursorPos {
			deletedItem = line
			curLine++
			continue
		}
		m.rw.WriteString(line)
		curLine++
	}
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	m.newAction("DEL", m.cursorPos, deletedItem)
}

func (m *Menu) CreateItem(offset int, content string) {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == m.cursorPos+offset {
			m.rw.WriteString(content)
		}
		m.rw.WriteString(line)
		curLine++
	}
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	m.newAction("NEW", m.cursorPos, "")
}

func (m *Menu) DrawMenu(isInteractive bool) (output string) {
	m.file.Seek(0, io.SeekStart)
	lineNum := 0
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if lineNum == m.cursorPos && isInteractive {
			output += "-> "
		} else {
			output += fmt.Sprintf(" %d ", lineNum+1)
		}
		if len(line) >= 5 && line[:5] == "[!-!]" {
			output += "[X] "
			line = line[5:]
		} else {
			output += "[ ] "
		}
		output += line
		lineNum++
	}
	if isInteractive {
		output += "Press 'h' for help\n"
	}
	return output
}

func (m *Menu) MarkItem() {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == m.cursorPos {
			if len(line) >= 5 && line[:5] == "[!-!]" {
				m.rw.WriteString(line[5:])
			} else {
				m.rw.WriteString("[!-!]" + line)
			}
		} else {
			m.rw.WriteString(line)
		}
		curLine++
	}
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	m.newAction("MARK", m.cursorPos, "")
}
