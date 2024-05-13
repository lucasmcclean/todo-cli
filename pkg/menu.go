package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Menu struct {
	File      os.File
	CursorPos int
	items     []*MenuItem
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
	newMenu.File = *file
	newMenu.CursorPos = 0

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		item := newMenuItem(scanner.Text(), i)
		newMenu.items = append(newMenu.items, item)
		i++
	}

	return newMenu
}

func (m *Menu) Draw() {
	fmt.Printf("\033[H\033[J")
	i := 0
	for ; i < m.CursorPos; i++ {
		fmt.Println("   [ ] " + m.items[i].text)
	}
	fmt.Printf(" > [ ] "+"\033[1m%s\033[0m\n", m.items[i].text)
	for i += 1; i < len(m.items); i++ {
		fmt.Println("   [ ] " + m.items[i].text)
	}
}

func (m *Menu) MoveCursor(delta int) {
	m.CursorPos += delta
	for m.CursorPos > len(m.items)-1 {
		m.CursorPos -= len(m.items)
	}
	for m.CursorPos < 0 {
		m.CursorPos += len(m.items)
	}
}

func (m *Menu) Close() error {
	err := m.File.Close()
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
