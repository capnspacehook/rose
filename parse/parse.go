package parse

//go: generate goyacc -o rose.y.go rose.y

import (
	"io"

	"github.com/capnspacehook/rose/ast"
)

var typeNames = map[string]bool{
	"any":    true,
	"int":    true,
	"float":  true,
	"string": true,
}

var boolConsts = map[string]bool{
	"false": false,
	"true":  true,
}

func Parse(in io.Reader) ([]ast.Statement, error) {
	lexer := NewLexer(in)
	yyParse(lexer)

	return lexer.Statements, lexer.Err()
}
