package comm

import (
	"fmt"
	"strings"
)

var MsgBus = NewMessageBus()

type MessageEvent struct {
	Source  string
	Target  string
	Action  string
	Message string
	ID      string
	Field   string
	Value   string
}

func (me MessageEvent) String() string {
	return fmt.Sprintf("[Src:%s -> Tar:%s] Action: %s, Msg: %s, ID: %s, Field: %s, Value: %s", me.Source, me.Target, me.Action, me.Message, me.ID, me.Field, me.Value)
}

func (me MessageEvent) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s", me.Source, me.Target, me.Action, me.Message, me.ID, me.Field, me.Value)
}

// Any object that wan't to listen for message needs to implement this method.
type IMessageListener interface {
	Listen(msg *MessageEvent)
}

type MessageBus struct {
	listeners []IMessageListener

	debug bool
}

func NewMessageBus() *MessageBus {
	mb := new(MessageBus)
	mb.listeners = []IMessageListener{}
	mb.debug = false
	return mb
}

func NewMessageFrom(csv string) *MessageEvent {
	me := new(MessageEvent)
	s := strings.Split(csv, ",")

	me.Source = s[0]
	me.Target = s[1]
	me.Action = s[2]
	me.Message = s[3]
	me.ID = s[4]
	me.Field = s[5]
	me.Value = s[6]

	return me
}

func (mb *MessageBus) Subscribe(listener IMessageListener) {
	mb.listeners = append(mb.listeners, listener)
}

func (mb *MessageBus) SendEvent(e MessageEvent) {
	me := new(MessageEvent)
	me.Source = e.Source
	me.Target = e.Target
	me.Action = e.Action
	me.Message = e.Message
	me.Field = e.Field
	me.Value = e.Value

	if mb.debug {
		fmt.Printf("SendEvent: %s\n", me)
	}

	for _, listener := range mb.listeners {
		listener.Listen(me)
	}
}

func (mb *MessageBus) Send(source, target, action, msg, value string) {
	me := new(MessageEvent)
	me.Source = source
	me.Target = target
	me.Action = action
	me.Message = msg
	me.Value = value
	if mb.debug {
		fmt.Printf("Send: %s\n", me)
	}

	for _, listener := range mb.listeners {
		listener.Listen(me)
	}
}

func (mb *MessageBus) Send2(source, target, action, msg, id, value string) {
	me := new(MessageEvent)
	me.Source = source
	me.Target = target
	me.Action = action
	me.Message = msg
	me.ID = id
	me.Value = value

	if mb.debug {
		fmt.Printf("Send2: %s\n", me)
	}
	for _, listener := range mb.listeners {
		listener.Listen(me)
	}
}

func (mb *MessageBus) Send3(source, target, action, msg, id, field, value string) {
	me := new(MessageEvent)
	me.Source = source
	me.Target = target
	me.Action = action
	me.Message = msg
	me.ID = id
	me.Field = field
	me.Value = value

	if mb.debug {
		fmt.Printf("Send3: %s\n", me)
	}

	for _, listener := range mb.listeners {
		listener.Listen(me)
	}
}

func (mb *MessageBus) Close() {
	mb.listeners = nil
}
