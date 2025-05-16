package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecripo/tisp-lang/lib/scanner"
)

func runPrompt() {
	lineScanner := bufio.NewScanner(os.Stdin)

	for {
		// fmt.Println("data: ")
		// scanner.Scan()
		// data := scanner.Text()
		fmt.Print("> ")
		if success := lineScanner.Scan(); !success {
			break
		}
		expr := lineScanner.Text()
		if expr == "" {
			continue
		} else if expr == "exit" || expr == "quit" {
			break
		}
		scanner := scanner.Scanner{Source: expr}
		tokens, err := scanner.Scan()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Println(tokens)
	}
}

func runFile(name string) {
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: tisp")
		os.Exit(64) // unix exit code for invalid number of arguments
	}
	if len(args) == 1 {
		runFile(args[1])
	} else {
		runPrompt()
	}

}
