package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/capnspacehook/rose/ast"
)

func printStatements(stmts []ast.Statement) {
	for _, stmt := range stmts {
		fmt.Println(stmt.String())
	}
}

func TestParser(t *testing.T) {
	input := `var foo int = 0b10011; var foo = 7;
var bar = 56.0e1
// ooh, underscores
var baz = 114_223_117
var assigned = foo
boolean = true
nothing = nil
char = '\x60'

const uwu = 5 // uwu
let yes = false
const goku float = 9000.01
let urMom bool = true

stringy = "hello there"
`

	statements, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	printStatements(statements)
}
