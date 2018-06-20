package main

import (
	"fmt"
	"os"

	"github.com/ei1chi/sample-lang/repl"
)

func main() {
	fmt.Printf("Hello! This is my language!\n")
	repl.Start(os.Stdin, os.Stdout)
}
