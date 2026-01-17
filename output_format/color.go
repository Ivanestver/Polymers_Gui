package output_format

import (
	"github.com/gookit/color"
)

const (
	resetCode  = "\033[0m"
	redCode    = "\033[31m"
	yellowCode = "\033[33m"
	purpleCode = "\033[35m"
)

func MakeRed(msg string) string {
	return color.Red.Sprint(msg)
}

func MakeYellow(msg string) string {
	return color.Yellow.Sprint(msg, yellowCode)
}

func MakePurple(msg string) string {
	return color.RGB(78, 0, 142).Sprint(msg)
}
