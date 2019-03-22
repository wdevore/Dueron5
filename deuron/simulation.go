package deuron

import (
	"github.com/wdevore/Deuron5/deuron/app/comm"
)

type ISimulation interface {
	Connect(statusChannel chan string)
	Command(args []string)
	Send(msg string)
	Create()
	Reset()
	Step()
	RunPause()
	SendEvent(event *comm.MessageEvent)
	ToJSON() interface{}
	Load(interface{})
}
