package main

// Note: You may need a ram drive:
// diskutil erasevolume HFS+ 'RAMDisk' `hdiutil attach -nomount ram://2097152`

import (
	"fmt"

	"github.com/wdevore/Deuron5/deuron/app"
)

/*
	Deuron has a GUI for display.
	Keyboard input for control
	Debug output to the console.

	typical sequence:
	>con
	>load
	>start
*/
var theApp *app.App // The GUI

// This is the main entry point for Deuron.
func main() {
	theApp = app.NewApp()
	defer theApp.Close()

	theApp.Open()

	theApp.SetFont("Roboto-Bold.ttf", 24)
	theApp.Configure()

	theApp.Run()
}

func printHelp() {
	fmt.Println("------------------- Help ---------------------------------")
	fmt.Println("'quit' stops any simulation and exits app.")
	fmt.Println("'help' this help screen.")
	fmt.Println("'start' starts the current target simulation.")
	fmt.Println("'stop' stops the current target simulation.")
	fmt.Println("'load' reads target's json.")
	fmt.Println("'go' connects, loads and runs sim.")
	fmt.Println("'set sim-name' sets the target simulation, where")
	fmt.Println("   sim-name specifies a json file in the working directory.")
	fmt.Println("'con' connects to a target sim-name. It does NOT start it.")
	fmt.Println("'type' changes sim type: `runreset` or `continous`")
	fmt.Println("'ping' sends `ping` to target sim.")

	// fmt.Println("'p' activates property mode and lists available properties.")
	// fmt.Println("  you then enter <property number> and <value>")
	fmt.Println("'\\' shows what properties can be changed. To change, for example,")
	fmt.Println("   Poisson-min value enter 1 1 <value>")
	fmt.Println("----------------------------------------------------------")
}
