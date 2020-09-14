package main

import (
	"fmt"
	"os"

	"github.com/capnspacehook/rose/repl"
)

func main() {
	fmt.Println("Rose REPL")
	repl.Start(os.Stdin, os.Stdout)
}
