package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Menu struct {
	file      os.File
	cursorPos int
	items     []*MenuItem
	completed []int
}

type MenuItem struct {
	text string
	pos  int
}

var file *os.File

func NewMenu(fileLocation string) *Menu {
	var err error
	file, err = os.Open(fileLocation)
	if err != nil {
		panic(errors.New("Could not open file " + fileLocation))
	}

	newMenu := &Menu{}
	newMenu.file = *file
	newMenu.cursorPos = 0
	newMenu.items = make([]*MenuItem, 0, 20)
	newMenu.completed = make([]int, 0, 20)

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		item := newMenuItem(scanner.Text(), i)
		newMenu.items = append(newMenu.items, item)
		i++
	}

	return newMenu
}

func (m *Menu) Draw(offset int) {
	if offset < 0 {
		offset = -offset
	}
	fmt.Printf("\033[%d;H\033[J", offset)
	for i := 0; i < len(m.items); i++ {
		completed := IncludesInt(m.completed, i)
		if i == m.cursorPos && completed {
			fmt.Printf(" > [X] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if i == m.cursorPos {
			fmt.Printf(" > [ ] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if completed {
			fmt.Println("   [X] " + m.items[i].text)
		} else {
			fmt.Println("   [ ] " + m.items[i].text)
		}
	}
}

func (m *Menu) MoveCursor(delta int) {
	if len(m.items) == 0 {
		m.cursorPos = 0
		return
	}
	m.cursorPos += delta
	for m.cursorPos > len(m.items)-1 {
		m.cursorPos -= len(m.items)
	}
	for m.cursorPos < 0 {
		m.cursorPos += len(m.items)
	}
}

func (m *Menu) CompleteItem(remove bool) {
	if remove {
		if len(m.items) <= 1 {
			m.items = m.items[:0]
		} else {
			m.items = append(m.items[:m.cursorPos], m.items[m.cursorPos+1:]...)
		}
		m.MoveCursor(0)
	} else {
		m.completed = append(m.completed, m.cursorPos)
	}
}

func (m *Menu) Close() error {
	err := m.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func newMenuItem(text string, pos int) *MenuItem {
	return &MenuItem{
		text: text,
		pos:  pos,
	}
}
