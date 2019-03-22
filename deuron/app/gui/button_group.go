package gui

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// Only one button can be selected at any one time.
type ButtonGroup struct {
	buttons *sll.List
}

func NewButtonGroup() *ButtonGroup {
	bg := new(ButtonGroup)

	bg.buttons = sll.New()

	return bg
}

func (bg *ButtonGroup) AddButton(button IWidget) {
	bg.buttons.Add(button)
}

func (bg *ButtonGroup) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	it := bg.buttons.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		handled, id = btn.Handle(x, y, eventType)
		if handled {
			bg.unselectAll()
			btn.Select()
			return handled, id
		}
	}

	return false, -1
}

func (bg *ButtonGroup) unselectAll() {
	it := bg.buttons.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.UnSelect()
	}
}

func (bg *ButtonGroup) Iterator() sll.Iterator {
	return bg.buttons.Iterator()
}
