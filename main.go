package main

import (
	"os"

	"github.com/scott-haines/sinject/commands"
)

// main entry point of the program
func main() {
	err := commands.NewSinjectCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
