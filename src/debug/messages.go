package debug

import "fmt"

func Message(msg string) {
	fmt.Println(msg)
}

func Success(msg string) {
	fmt.Println(Colorize(GREEN, "SUCCESS: "), msg)
}

func ErrorMessage(msg string) {
	fmt.Println(Colorize(RED, "ERROR: "), msg)
}

func Error(msg error) {
	fmt.Println(Colorize(RED, "ERROR: "), msg)
}
