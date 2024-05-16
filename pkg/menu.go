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
	undoStack []*MenuItemHist
	limbo     *MenuItem
}

type MenuItem struct {
	text      string
	completed bool
}

type MenuItemHist struct {
	menuItem  *MenuItem
	pos       int
	operation string
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
		if i == m.cursorPos && m.items[i].completed {
			fmt.Printf(" > [X] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if i == m.cursorPos {
			fmt.Printf(" > [ ] "+"\033[1m%s\033[0m\n", m.items[i].text)
		} else if m.items[i].completed {
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
		m.postAction(m.cursorPos, m.items[m.cursorPos], "DEL")
		if len(m.items) <= 1 {
			m.items = m.items[:0]
			return
		}
		m.items = append(m.items[:m.cursorPos], m.items[m.cursorPos+1:]...)
		m.MoveCursor(0)
	} else {
		m.items[m.cursorPos].completed = true
	}
}

func (m *Menu) AddItem(rPos int) {
	if m.limbo != nil {
		m.placeItem(rPos)
		return
	}
	input := GetInput()
	if input == "" {
		return
	}
	newItem := newMenuItem(input)
	if len(m.items) == 0 {
		m.items = append(m.items, newItem)
		return
	}
	remainder := []*MenuItem{newItem}
	remainder = append(remainder, m.items[m.cursorPos+rPos:]...)
	m.items = append(m.items[:m.cursorPos+rPos], remainder...)
	m.postAction(m.cursorPos+rPos, newItem, "ADD")
}

func (m *Menu) MoveItem() {
	if m.limbo != nil {
		return
	}
	m.limbo = m.items[m.cursorPos]
	m.items = append(m.items[:m.cursorPos], m.items[m.cursorPos+1:]...)
}

func (m *Menu) placeItem(rPos int) {
	slice := []*MenuItem{m.limbo}
	slice = append(slice, m.items[m.cursorPos+rPos:]...)
	m.items = append(m.items[:m.cursorPos+rPos], slice...)
	m.postAction(m.cursorPos, m.limbo, "MOV")
	m.limbo = nil
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

func (m *Menu) postAction(pos int, item *MenuItem, op string) {
	hist := &MenuItemHist{
		menuItem:  item,
		pos:       pos,
		operation: op,
	}
	m.undoStack = append(m.undoStack, hist)
}

func (m *Menu) UndoAction() {
	var action *MenuItemHist
	action, m.undoStack = m.undoStack[len(m.undoStack)-1], m.undoStack[:len(m.undoStack)-1]
	switch action.operation {
	case "DEL":
		slice := []*MenuItem{action.menuItem}
		slice = append(slice, m.items[action.pos:]...)
		m.items = append(m.items[:action.pos], slice...)
	case "ADD":
		m.items = append(m.items[:action.pos], m.items[action.pos+1:]...)
	case "MOV":
		m.items = append(m.items[:action.pos], m.items[action.pos+1:]...)
	}
}
