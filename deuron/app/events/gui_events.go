package events

type MouseEventType int

const (
	MouseMotion MouseEventType = iota
	MouseButton
	MouseWheel
)
