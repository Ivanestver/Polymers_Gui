package output_format

type IPrint interface {
	Print(str string)
	Println(str string)
	Printf(format string, args ...any)
	Printfln(format string, args ...any)
	PrintEmptyLine()
	PrintlnInfo(msg string)
	PrintInfo(msg string)
	PrintfInfo(msg string, args ...any)
	PrintlnError(msg string)
	PrintflnError(msg string, args ...any)
	PrintlnWarning(msg string)
	PrintWarning(msg string)
	PrintfWarning(msg string, args ...any)
	PrintflnWarning(msg string, args ...any)
	Read(args ...any)
	Readln(args ...any)
}
