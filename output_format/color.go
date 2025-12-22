package output_format

const (
	resetCode  = "\033[0m"
	redCode    = "\033[31m"
	yellowCode = "\033[33m"
	purpleCode = "\033[35m"
)

func MakeColored(msg, colorCode string) string {
	return colorCode + msg + resetCode
}

func MakeRed(msg string) string {
	return MakeColored(msg, redCode)
}

func MakeYellow(msg string) string {
	return MakeColored(msg, yellowCode)
}

func MakePurple(msg string) string {
	return MakeColored(msg, purpleCode)
}
