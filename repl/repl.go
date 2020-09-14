package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/capnspacehook/rose/lexer"
	"github.com/capnspacehook/rose/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var l lexer.Lexer
	fs := token.NewFileSet()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l.Init(fs.AddFile("", fs.Base(), len(line)), strings.NewReader(line), true)

		for tok, lit := l.Lex(); tok != token.EOF; tok, lit = l.Lex() {
			if err := l.Err(); err != nil {
				fmt.Fprintf(out, "error: %v\n", err)
				continue
			}

			fmt.Fprintf(out, `"%+v": %q`+"\n", tok, lit)
		}
	}
}
