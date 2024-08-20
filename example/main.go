package main

import (
	"fmt"

	"github.com/CornMars2020/keyboard"
)

func cmdHandler(cmd string) {
	fmt.Println(cmd)

	// customize your own command executer
}

func helpFunc() {
	fmt.Println("This is Help")

	// customize your own help message
}

func main() {
	// set cmd alias
	keyboard.SetFastCmd("sws", "set workspace settings")
	keyboard.SetFastCmd("pws", "print workspace settings")

	// set help func
	keyboard.SetHelpFunc(helpFunc)
	// start program
	keyboard.HandleKeyboard(cmdHandler)
}
