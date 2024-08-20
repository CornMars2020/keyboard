package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CornMars2020/color"
	"github.com/CornMars2020/keyboard"
)

func setCmd(cmds []string) {
	fmt.Println("set cmd:", cmds)
}

func getCmd(cmds []string) {
	fmt.Println("get cmd:", cmds)
}

func cmdHandler(cmd string) {
	fmt.Println(cmd)

	// customize your own command executer
	cmds := strings.Split(cmd, " ")

	switch cmds[0] {

	case "set":
		setCmd(cmds[1:])
	case "get":
		getCmd(cmds[1:])

	default:
		log.Printf(color.GetRed("unknown cmd: %s"), cmd)
	}
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
