package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/capnspacehook/rose/token"
)

func TestLexer(t *testing.T) {
	input := `foo = 8
nine = "hi"
	`
	fs := token.NewFileSet()

	var l lexer
	l.Init(fs.AddFile("", fs.Base(), len(input)), strings.NewReader(input), len(input))
	for i := 0; i < 9; i++ {
		pos, tok, lit := l.Lex()
		fmt.Println(pos, tok, lit)
	}

	err := l.Err()
	if err != nil {
		t.Error(err)
	}
}
