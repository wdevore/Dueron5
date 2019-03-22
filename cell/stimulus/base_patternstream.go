package stimulus

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron5/cell"
)

type basePatternStream struct {
	id    int
	cons  *sll.List //[]cell.IConnection
	value int
}

func (ba *basePatternStream) baseInitialize() {
	ba.cons = sll.New()
}

func (ba *basePatternStream) Attach(con cell.IConnection) {
	ba.cons.Add(con)
}

func (ba *basePatternStream) SetId(id int) {
	ba.id = id
}

func (ba *basePatternStream) Id() int {
	return ba.id
}
