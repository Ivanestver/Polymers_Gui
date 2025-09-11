package output_format

import "fmt"

func PrintlnInfo(msg string) {
	fmt.Println("[INFO] " + msg)
}

func PrintInfo(msg string) {
	fmt.Print("[INFO] " + msg)
}

func PrintfInfo(msg string, args ...any) {
	fmt.Printf("[INFO] "+msg, args)
}

func PrintlnError(msg string) {
	fmt.Println("[ERROR] " + msg)
}

func PrintflnError(msg string, args ...any) {
	fmt.Printf("[ERROR] "+msg+"\n", args)
}

func PrintEmptyLine() {
	fmt.Println("")
}
