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
		history:   []action{},
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

func (m *Menu) DeleteItem(pos int, undo bool) {
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
	if !undo && curLine != 0 {
		m.newAction("DEL", func() { m.CreateItem(pos, deletedItem, true) })
	}
}

func (m *Menu) CreateItem(pos int, content string, undo bool) {
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
	if curLine <= pos {
		m.rw.WriteString(content)
	}
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	if !undo {
		m.newAction("NEW", func() { m.DeleteItem(pos, true) })
	}
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
			output += " -> "
		} else {
			if lineNum < 9 {
				output += fmt.Sprintf("  %d ", lineNum+1)
			} else {
				output += fmt.Sprintf(" %d ", lineNum+1)
			}
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

func (m *Menu) MarkItem(pos int, undo bool) {
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
	if !undo {
		m.newAction("MARK", func() { m.MarkItem(pos, true) })
	}
}

func (m *Menu) MoveItem(oldPos int, newPos int, undo bool) {
	m.file.Seek(0, io.SeekStart)
	curLine := 0
	var numBefore int
	var before, moved, after string
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			break
		}
		if curLine == oldPos {
			moved = line
		} else if curLine < newPos {
			before += line
			numBefore++
		} else if curLine >= newPos {
			after += line
		}
		curLine++
	}
	m.rw.WriteString(before + moved + after)
	m.file.Truncate(0)
	m.file.Seek(0, io.SeekStart)
	m.rw.Flush()
	if !undo {
		if newPos < oldPos {
			oldPos++
		}
		m.newAction("MOV", func() { m.MoveItem(numBefore, oldPos, true) })
	}
}

func (m *Menu) Undo() {
	action, err := m.undoAction()
	if err != nil {
		return
	}
	action.undo()
}

func (m *Menu) undoAction() (act action, err error) {
	if len(m.history) == 0 {
		return act, fmt.Errorf("There is nothing to undo")
	}
	act = m.history[len(m.history)-1]
	m.history = m.history[:len(m.history)-1]
	return act, nil
}
