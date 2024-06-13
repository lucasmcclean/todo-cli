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
	op   string
	undo func()
}

func (m *Menu) newAction(op string, undo func()) {
	action := action{
		op:   op,
		undo: undo,
	}
	m.history = append(m.history, action)
}

func (m *Menu) undoAction() (action action) {
	action = m.history[len(m.history)-1]
	m.history = append(m.history, m.history[:len(m.history)-1]...)
	return action
}

func (m *Menu) MoveCursor(delta int) {
	m.cursorPos += delta
	fileLen := GetFileLength(*m.file)
	for m.cursorPos < 0 {
		m.cursorPos += fileLen
	}
	for m.cursorPos >= fileLen {
		m.cursorPos -= fileLen
	}
}

func (m *Menu) DeleteItem(pos int) {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	var deletedItem string
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == pos {
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
	m.newAction("DEL", func() { m.CreateItem(m.cursorPos, deletedItem) })
}

func (m *Menu) CreateItem(pos int, content string) {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == pos {
			m.rw.WriteString(content)
		}
		m.rw.WriteString(line)
		curLine++
	}
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	m.newAction("NEW", func() { m.DeleteItem(pos) })
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

func (m *Menu) MarkItem(pos int) {
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
	m.newAction("MARK", func() { m.MarkItem(pos) })
}

func (m *Menu) MoveItem(oldPos int, newPos int) {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	var before, after, old string
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == oldPos {
			old = line
		} else if curLine < newPos {
			before += line
		} else if curLine >= newPos {
			after += line
		}
		curLine++
	}
	m.rw.WriteString(before + old + after)
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	m.newAction("MOV", func() { m.MoveItem(newPos, oldPos) })
}

func (m *Menu) Undo() {
	action := m.undoAction()
	action.undo()
}
