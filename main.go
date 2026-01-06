package main

import (
	"os"
	"typc/repl"
)

func main() {
	repl.Run(os.Stdin, os.Stdout)
}
