package output_format

import "fmt"

type ColoredConsolePrint struct{}

func (print *ColoredConsolePrint) Print(str string) {
	fmt.Print(str)
}

func (print *ColoredConsolePrint) Println(str string) {
	fmt.Println(str)
}

func (print *ColoredConsolePrint) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (print *ColoredConsolePrint) Printfln(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func (print *ColoredConsolePrint) PrintEmptyLine() {
	print.Println("")
}

func (print *ColoredConsolePrint) PrintlnInfo(msg string) {
	print.Println(MakeYellow("[INFO] ") + msg)
}

func (print *ColoredConsolePrint) PrintInfo(msg string) {
	print.Print(MakeYellow("[INFO] ") + msg)
}

func (print *ColoredConsolePrint) PrintfInfo(msg string, args ...any) {
	print.Printf(MakeYellow("[INFO] ")+msg, args)
}

func (print *ColoredConsolePrint) PrintlnError(msg string) {
	print.Println(MakeRed("[ERROR] " + msg))
}

func (print *ColoredConsolePrint) PrintflnError(msg string, args ...any) {
	print.Printfln(MakeRed("[ERROR] "+msg), args)
}

func (print *ColoredConsolePrint) Read(args ...any) {
	fmt.Scan(args...)
}

func (print *ColoredConsolePrint) Readln(args ...any) {
	fmt.Scanln(args...)
}
