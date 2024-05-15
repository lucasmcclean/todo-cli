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
		item := newMenuItem(scanner.Text())
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
	fmt.Println("  j: down, k: up, x: remove, enter: check,\n  i: new item before, a: new item after\n--------------------")
	for i := 0; i < len(m.items); i++ {
		isCompleted := IncludesInt(m.completed, i)
		if i == m.cursorPos && isCompleted {
			fmt.Printf(" > [X] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if i == m.cursorPos {
			fmt.Printf(" > [ ] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if isCompleted {
			fmt.Println("   [X] " + m.items[i].text)
		} else {
			fmt.Println("   [ ] " + m.items[i].text)
		}
	}
	fmt.Println("--------------------\n  'q' to quit")
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
		RemoveInt(&m.completed, m.cursorPos)
		if len(m.items) <= 1 {
			m.items = m.items[:0]
			return
		}
		m.items = append(m.items[:m.cursorPos], m.items[m.cursorPos+1:]...)
		m.MoveCursor(0)
	} else {
		if IncludesInt(m.completed, m.cursorPos) {
			RemoveInt(&m.completed, m.cursorPos)
			return
		}
		m.completed = append(m.completed, m.cursorPos)
	}
}

func (m *Menu) AddItem(rPos int, text string) {
	if text == "" {
		return
	}
	newItem := newMenuItem(text)
	if len(m.items) == 0 {
		m.items = append(m.items, newItem)
		return
	}
	// TODO: Cleaner solution?
	remainder := []*MenuItem{newItem}
	remainder = append(remainder, m.items[m.cursorPos+rPos:]...)
	m.items = append(m.items[:m.cursorPos+rPos], remainder...)
	for i := 0; i < len(m.completed); i++ {
		if m.completed[i] > m.cursorPos+rPos-1 {
			m.completed[i] += 1
		}
	}
}

func (m *Menu) Close() error {
	err := m.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func newMenuItem(text string) *MenuItem {
	return &MenuItem{
		text: text,
	}
}

func GetInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("New item: ")
	scanner.Scan()
	return scanner.Text()
}
