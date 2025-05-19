package main

import "fmt"

func printlnInfo(msg string) {
	fmt.Println("[INFO] " + msg)
}

func printInfo(msg string) {
	fmt.Print("[INFO] " + msg)
}

func printfInfo(msg string, args ...any) {
	fmt.Printf("[INFO] "+msg, args)
}

func printlnError(msg string) {
	fmt.Println("[ERROR] " + msg)
}

func printflnError(msg string, args ...any) {
	fmt.Printf("[ERROR] "+msg+"\n", args)
}

func printEmptyLine() {
	fmt.Println("")
}
