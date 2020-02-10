package parse

//go:generate goyacc -o rose.y.go rose.y

import (
	"io"

	"github.com/capnspacehook/rose/ast"
)

var typeNames = map[string]bool{
	"any":    true,
	"bool":   true,
	"int":    true,
	"float":  true,
	"string": true,
}

var boolConsts = map[string]bool{
	"false": false,
	"true":  true,
}

func Parse(in io.Reader) (*ast.Program, error) {
	lexer := newLexer(in)
	yyParse(lexer)

	return lexer.Program, lexer.Err()
}
