package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/capnspacehook/rose/parser"
	"github.com/capnspacehook/rose/token"

	"github.com/capnspacehook/pretty"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		fset := token.NewFileSet()
		file := fset.AddFile("", -1, len(line))
		ast, err := parser.ParseFile(file, strings.NewReader(line))
		if err != nil {
			fmt.Fprintf(out, "error: %v\n", err)
			continue
		}

		pretty.Println(ast)
	}
}
