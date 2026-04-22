package outputformat

import (
	"fmt"
	"os"
)

type FilePrint struct {
	file *os.File
}

func NewFilePrint(file *os.File) *FilePrint {
	return &FilePrint{
		file: file,
	}
}

func (print *FilePrint) Print(str string) {
	fmt.Fprint(print.file, str)
}

func (print *FilePrint) Println(str string) {
	fmt.Fprintln(print.file, str)
}

func (print *FilePrint) Printf(format string, args ...any) {
	fmt.Fprintf(print.file, format, args...)
}

func (print *FilePrint) Printfln(format string, args ...any) {
	fmt.Fprintf(print.file, format+"\n", args...)
}

func (print *FilePrint) PrintEmptyLine() {
	print.Println("")
}

func (print *FilePrint) PrintlnInfo(msg string) {
	print.Println("[INFO] " + msg)
}

func (print *FilePrint) PrintInfo(msg string) {
	print.Print("[INFO] " + msg)
}

func (print *FilePrint) PrintfInfo(msg string, args ...any) {
	print.Printf("[INFO] "+msg, args...)
}

func (print *FilePrint) PrintflnInfo(msg string, args ...any) {
	print.PrintfInfo("[INFO] "+msg+"\n", args...)
}

func (print *FilePrint) PrintlnError(msg string) {
	print.Println("[ERROR] " + msg)
}

func (print *FilePrint) PrintflnError(msg string, args ...any) {
	print.Printfln("[ERROR] "+msg, args...)
}

func (print *FilePrint) PrintlnWarning(msg string) {
	print.Println("[WARNING] " + msg)
}

func (print *FilePrint) PrintWarning(msg string) {
	print.Print("[WARNING] " + msg)
}

func (print *FilePrint) PrintfWarning(msg string, args ...any) {
	print.Printf("[WARNING] "+msg, args...)
}

func (print *FilePrint) PrintflnWarning(msg string, args ...any) {
	print.Printfln("[WARNING] "+msg, args...)
}

func (print *FilePrint) Read(args ...any) {
	if _, err := fmt.Fscan(print.file, args); err != nil {
		panic(err)
	}
}

func (print *FilePrint) Readln(args ...any) {
	if _, err := fmt.Fscanln(print.file, args); err != nil {
		panic(err)
	}
}
