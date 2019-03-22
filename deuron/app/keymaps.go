package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (ap *App) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		fmt.Println("SDL Quit event")
		ap.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		ap.mouseX = t.X
		ap.mouseY = t.Y
		ap.Handle(t.X, t.Y, events.MouseMotion)

		// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseButtonEvent:
		if t.State == 1 {
			ap.Handle(t.X, t.Y, events.MouseButton)
			// fmt.Printf("[%d ms](%v) MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
			// 	t.Timestamp, handled, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		}
		return false
	case *sdl.MouseWheelEvent:
		ap.Handle(t.X, t.Y, events.MouseWheel)
		// fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y)
		return false
	case *sdl.KeyboardEvent:
		// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tScode:%d\tstate:%d\trepeat:%d\n",
		// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.Keysym.Scancode, t.State, t.Repeat)
		switch t.Keysym.Scancode {
		case sdl.SCANCODE_LEFT:
			sync := deuron.SimModel.GetFloat("RangeSync")

			if sync == 1 {
				ap.rangeScrollLeft()
			} else {
				// Scroll graph left. The graphs listen for this type of message.
				comm.MsgBus.Send2("Keymaps", "Key", "Left", "ScrollLeft", ap.controlID, "")
			}

			return false
		case sdl.SCANCODE_RIGHT:
			sync := deuron.SimModel.GetFloat("RangeSync")

			if sync == 1 {
				ap.rangeScrollRight()
			} else {
				// Scroll graph right
				comm.MsgBus.Send2("Keymaps", "Key", "Right", "ScrollRight", ap.controlID, "")
			}

			return false
		case sdl.SCANCODE_R:
			if t.State == sdl.RELEASED {
				comm.MsgBus.Send("Keymaps", "App", "Command", "RunPause", "")
			}
			return false
		}

		if t.State == sdl.RELEASED {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				// Canceling edit mode on a value button.
				ap.currentValue = ""
				comm.MsgBus.Send2("Keymaps", "Control", "Edit", "Cancel", ap.controlID, ap.controlOrgValue)
				break
			case sdl.SCANCODE_GRAVE: // The tilte or single quote key
				fmt.Println("KeyEscape Quit event")
				ap.running = false
				break
			case sdl.SCANCODE_RETURN: // Signals finished entering data
				// Update control to change modes.
				comm.MsgBus.Send3("Keymaps", "Control", "Edit", "Done", ap.controlID, ap.controlField, ap.currentValue)
				// Update model. Model will relay to listeners about change.
				comm.MsgBus.Send3("Keymaps", "Model", "Set", ap.controlMsg, ap.controlID, ap.controlField, ap.currentValue)
				break
			default:
				ap.currentValue = ap.currentValue + codeToString(t.Keysym.Scancode)
				comm.MsgBus.Send3("Keymaps", "Control", "Edit", "Current", ap.controlID, ap.controlField, ap.currentValue)
				break
			}
		}

		return false // We handled it. Don't allow it to be added to the queue.
	}

	return true
}

func (ap *App) rangeScrollLeft() {
	start := int(deuron.SimModel.GetFloat("Range_Start"))
	end := int(deuron.SimModel.GetFloat("Range_End"))

	if start > 0 {
		inc := deuron.SimModel.GetFloat("Inc/Dec")

		start = start - int(inc)
		if start >= 0 {
			end = end - int(inc)
			comm.MsgBus.Send3("Keymaps", "Model", "Set", "", "", "Range_Start", fmt.Sprintf("%d", start))
			comm.MsgBus.Send3("Keymaps", "Model", "Set", "", "", "Range_End", fmt.Sprintf("%d", end))
		}
	}
}

func (ap *App) rangeScrollRight() {
	sync := deuron.SimModel.GetFloat("RangeSync")

	if sync == 0 {
		return
	}

	start := int(deuron.SimModel.GetFloat("Range_Start"))
	end := int(deuron.SimModel.GetFloat("Range_End"))
	duration := int(deuron.SimModel.GetFloat("Samples"))

	if end <= duration { //  or <
		inc := deuron.SimModel.GetFloat("Inc/Dec")
		end = end + int(inc)
		if end <= duration {
			start = start + int(inc)
			comm.MsgBus.Send3("Keymaps", "Model", "Set", "", "", "Range_Start", fmt.Sprintf("%d", start))
			comm.MsgBus.Send3("Keymaps", "Model", "Set", "", "", "Range_End", fmt.Sprintf("%d", end))
		}
	}
}

func codeToInt(code sdl.Scancode) int {
	switch code {
	case sdl.SCANCODE_0:
		return 0
	case sdl.SCANCODE_1:
		return 1
	case sdl.SCANCODE_2:
		return 2
	case sdl.SCANCODE_3:
		return 3
	case sdl.SCANCODE_4:
		return 4
	case sdl.SCANCODE_5:
		return 5
	case sdl.SCANCODE_6:
		return 6
	case sdl.SCANCODE_7:
		return 7
	case sdl.SCANCODE_8:
		return 8
	case sdl.SCANCODE_9:
		return 9
	}
	return -1
}

func codeToString(code sdl.Scancode) string {
	switch code {
	case sdl.SCANCODE_0:
		return "0"
	case sdl.SCANCODE_1:
		return "1"
	case sdl.SCANCODE_2:
		return "2"
	case sdl.SCANCODE_3:
		return "3"
	case sdl.SCANCODE_4:
		return "4"
	case sdl.SCANCODE_5:
		return "5"
	case sdl.SCANCODE_6:
		return "6"
	case sdl.SCANCODE_7:
		return "7"
	case sdl.SCANCODE_8:
		return "8"
	case sdl.SCANCODE_9:
		return "9"
	case sdl.SCANCODE_PERIOD:
		return "."
	}
	return ""
}
