package output_format

var iprint IPrint

func SetPrint(p IPrint) {
	iprint = p
}

func Print(str string) {
	iprint.Print(str)
}

func Println(str string) {
	iprint.Println(str)
}

func Printf(format string, args ...any) {
	iprint.Printf(format, args...)
}

func Printfln(format string, args ...any) {
	iprint.Printfln(format, args...)
}

func PrintEmptyLine() {
	iprint.PrintEmptyLine()
}

func PrintlnInfo(msg string) {
	iprint.PrintlnInfo(msg)
}

func PrintInfo(msg string) {
	iprint.PrintInfo(msg)
}

func PrintfInfo(msg string, args ...any) {
	iprint.PrintfInfo(msg, args)
}

func PrintlnError(msg string) {
	iprint.PrintlnError(msg)
}

func PrintflnError(msg string, args ...any) {
	iprint.PrintflnError(msg, args)
}

func Read(args ...any) {
	iprint.Read(args...)
}

func Readln(args ...any) {
	iprint.Readln(args...)
}
