package debug

type Color int

const (
	RESET Color = 0
	RED   Color = 1
	GREEN Color = 2
)

const reset_s string = "\033[0m"

const red_s string = "\033[31m"
const green_s string = "\033[32m"

func Colorize(c Color, msg string) string {
	if c == RED {
		return red_s + msg + reset_s
	}
	if c == GREEN {
		return green_s + msg + reset_s
	}
	return ""
}
