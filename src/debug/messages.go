package debug

import "fmt"

func PrintMsg(msg string) {
	fmt.Println(msg)
}

func PrintSuccess(msg string) {
	fmt.Println(Colorize(GREEN, "SUCCESS: "), msg)
}

func PrintErrorMsg(msg string) {
	fmt.Println(Colorize(RED, "ERROR: "), msg)
}

func PrintError(msg error) {
	fmt.Println(Colorize(RED, "ERROR: "), msg)
}
