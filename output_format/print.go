package output_format

import "fmt"

func Print(str string) {
	fmt.Print(str)
}

func Println(str string) {
	fmt.Println(str)
}

func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func Printfln(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func PrintEmptyLine() {
	Println("")
}

func PrintlnInfo(msg string) {
	Println(MakeYellow("[INFO] ") + msg)
}

func PrintInfo(msg string) {
	Print(MakeYellow("[INFO] ") + msg)
}

func PrintfInfo(msg string, args ...any) {
	Printf(MakeYellow("[INFO] ")+msg, args)
}

func PrintlnError(msg string) {
	Println(MakeRed("[ERROR] " + msg))
}

func PrintflnError(msg string, args ...any) {
	Printfln(MakeRed("[ERROR] "+msg), args)
}

func Read(args ...any) {
	fmt.Scan(args...)
}

func Readln(args ...any) {
	fmt.Scanln(args...)
}
