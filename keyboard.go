package keyboard

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CornMars2020/color"
	"github.com/eiannone/keyboard"
)

var str string
var commandSlice []string = make([]string, 0)
var commandIndex int = 0
var helpCmdFunc func() = func() {}

var historyFileDir string = "./"

// fast command alias
var fastCmd map[string]string = map[string]string{
	"ls": "help",
	"h":  "help",
	"q":  "quit",
}

// not record into history
var noLogHistoryCmd []string = []string{
	"clear",
	"q", "quit", "exit",
	"ls", "h", "help",
}

// valid chars for input
var validChars []string = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "=", "+",
	"[", "]", "{", "}", "|", "\\", ":", ";", "'", "\"", "<", ">", ",", ".", "?", "/",
}

func init() {
	// recover from file
	stat, err := os.Stat(historyFileDir)
	if err == nil && stat.IsDir() {
		log.Println(color.GetYellow("load from .history"))
		loadCommand()
		return
	}

	historyFileDir = "." + historyFileDir
	stat, err = os.Stat(historyFileDir)
	if err == nil && stat.IsDir() {
		log.Println(color.GetYellow("load from .history"))
		loadCommand()
		return
	}

	execFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	historyFileDir = path.Join(filepath.Dir(execFile), "../") + "/"
	stat, err = os.Stat(historyFileDir)
	if err == nil && stat.IsDir() {
		log.Println(color.GetYellow("load from .history"))
		loadCommand()
		return
	}

	log.Println(color.GetRed("no .history file"))
}

func isValidChar(str string) bool {
	if len(str) != 1 {
		return false
	}

	for _, v := range validChars {
		if v == str {
			return true
		}
	}

	return false
}

func loadCommand() {
	histFile := historyFileDir + ".history"
	content, err := os.ReadFile(histFile)
	if err != nil {
		log.Println(color.GetRed(err.Error()))
		return
	}

	commands := strings.Split(string(content), "\n")
	for _, v := range commands {
		if v != "" {
			commandSlice = append(commandSlice, v)
		}
	}
}

func saveCommand(str string) {
	for _, v := range noLogHistoryCmd {
		if str == v {
			return
		}
	}

	commandSlice = append(commandSlice, str)
	if len(commandSlice) > 100 {
		commandSlice = (commandSlice)[1:]
	}

	histFile := historyFileDir + ".history"
	content := strings.Join(commandSlice, "\n")
	err := os.WriteFile(histFile, []byte(content), 0644)
	if err != nil {
		log.Println(color.GetRed(err.Error()))
		return
	}
}

func getCommand(up bool) (commamd string) {
	if up {
		if commandIndex > 0 {
			commandIndex--
		}
	} else {
		if commandIndex < len(commandSlice) {
			commandIndex++
		}
	}

	if commandIndex < len(commandSlice) {
		return commandSlice[commandIndex]
	} else {
		return ""
	}
}

func handleCmd(cmd string, executeCmd func(string)) string {
	// 处理命令
	if len(cmd) <= 0 {
		return "continue"
	}
	inputCmd := cmd[:]

	cmd = strings.TrimSpace(inputCmd)
	re, _ := regexp.Compile("[ \t]+")
	cmd = re.ReplaceAllString(cmd, " ")

	if inputCmd != cmd {
		log.Printf(color.GetCyan("cmd: '%s' => '%s'"), inputCmd, cmd)
	}

	// 快捷命令 别名
	if fCmd, ok := fastCmd[cmd]; ok {
		cmd = fCmd
	}

	// 基础命令: clear, help, quit/exit

	if cmd == "clear" {
		fmt.Print("\033[H\033[2J")
		return "continue"
	}

	if cmd == "help" {
		if helpCmdFunc != nil {
			helpCmdFunc()
		}

		log.Println("退出: \tq|quit|exit|ESC")
		return "continue"
	}

	if strings.Contains(cmd, "quit") || strings.Contains(cmd, "exit") {
		log.Printf(color.GetRed("%s: process exiting"), cmd)
		return "break"
	}

	// 复杂命令
	executeCmd(cmd)
	return "continue"
}

func SetFastCmd(cmd string, aliasCmd string) {
	cmd = strings.ToLower(cmd)
	aliasCmd = strings.ToLower(aliasCmd)
	fastCmd[cmd] = aliasCmd
}

func SetHelpFunc(help func()) {
	helpCmdFunc = help
}

func HandleKeyboard(executeCmd func(str string)) {
	// 交互式命令

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	updateTm := time.Now()
	go func() {
		for {
			if time.Since(updateTm) > 300*time.Second {
				_ = keyboard.Close()

				os.Exit(0)
			}
			time.Sleep(time.Second)
		}
	}()

	fmt.Print("> ")

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		updateTm = time.Now()

		// esc 退出
		if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC || event.Key == keyboard.KeyCtrlD {
			break
		}

		// 上下翻页
		if event.Key == keyboard.KeyArrowUp {
			command := getCommand(true)
			str = command
			fmt.Print("\033[2K\r", "> ", command)
		}
		if event.Key == keyboard.KeyArrowDown {
			command := getCommand(false)
			str = command
			fmt.Print("\033[2K\r", "> ", command)
		}

		// 左右光标移动 TODO: 移动后输入, 显示覆盖了后面字符, 实际是字符拼接
		// if event.Key == keyboard.KeyArrowLeft {
		// 	fmt.Print("\033[1D")
		// }
		// if event.Key == keyboard.KeyArrowRight {
		// 	fmt.Print("\033[1C")
		// }

		// 输入正常字符
		if isValidChar(string(event.Rune)) {
			fmt.Print(string(event.Rune))
			str = str + string(event.Rune)
		}

		// 输入空格
		if event.Key == keyboard.KeySpace {
			fmt.Print(" ")
			str = str + " "
		}

		// 删除
		if event.Key == keyboard.KeyBackspace || event.Key == keyboard.KeyBackspace2 {
			if len(str) >= 1 {
				str = str[:len(str)-1]
			} else {
				str = ""
			}
			fmt.Print("\033[2K\r", "> ", str)
		}

		// 回车
		if event.Key == keyboard.KeyEnter {
			fmt.Printf("\r\n")

			if str != "" {
				saveCommand(str)
				commandIndex = len(commandSlice)

				flag := handleCmd(str, executeCmd)
				if flag == "break" {
					break
				}
			}

			str = ""
			fmt.Print("> ")
		}
	}
}
