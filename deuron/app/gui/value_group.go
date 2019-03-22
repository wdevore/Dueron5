package gui

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// Collects Edit/Value widgets
type ValueGroup struct {
	values *sll.List
}

func NewValueGroup() *ValueGroup {
	bg := new(ValueGroup)

	bg.values = sll.New()

	return bg
}

func (bg *ValueGroup) AddWidget(widget IWidget) {
	bg.values.Add(widget)
}

func (bg *ValueGroup) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	it := bg.values.Iterator()
	for it.Next() {
		btn := it.Value().(*ValueButton)
		handled, id = btn.Handle(x, y, eventType)
		if handled {
			return handled, id
		}
	}

	return false, -1
}

func (bg *ValueGroup) Iterator() sll.Iterator {
	return bg.values.Iterator()
}
