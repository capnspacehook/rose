package lexer

import (
	"strings"
	"testing"

	"github.com/capnspacehook/rose/token"
)

func TestLex(t *testing.T) {
	input := `t = true
const c = "constant";
a, b = nil, 99.53

fn arithmetic(x, y int) int {
	return x + y
}

if 5 < 9 {
	return true
} else {
	return false
}

[1, 2, 3]
{7, 7, 8}
{1: "one", 2: "two"};
("tru", 8)
l[2:5]
`

	tests := []struct {
		tok token.Token
		lit string
	}{
		{token.IDENT, "t"},
		{token.ASSIGN, ""},
		{token.IDENT, "true"},
		{token.SEMI, "\n"},
		{token.CONST, "const"},
		{token.IDENT, "c"},
		{token.ASSIGN, ""},
		{token.STRING, `"constant"`},
		{token.SEMI, ";"},
		{token.IDENT, "a"},
		{token.COMMA, ""},
		{token.IDENT, "b"},
		{token.ASSIGN, ""},
		{token.IDENT, "nil"},
		{token.COMMA, ""},
		{token.FLOAT, "99.53"},
		{token.SEMI, "\n"},

		{token.FN, "fn"},
		{token.IDENT, "arithmetic"},
		{token.LPAREN, ""},
		{token.IDENT, "x"},
		{token.COMMA, ""},
		{token.IDENT, "y"},
		{token.IDENT, "int"},
		{token.RPAREN, ""},
		{token.IDENT, "int"},
		{token.LBRACE, ""},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.ADD, ""},
		{token.IDENT, "y"},
		{token.SEMI, "\n"},
		{token.RBRACE, ""},
		{token.SEMI, "\n"},
	}

	var l Lexer
	fs := token.NewFileSet()

	l.Init(fs.AddFile("", fs.Base(), len(input)), strings.NewReader(input), true)
	for _, test := range tests {
		tok, lit := l.Lex()
		if tok != test.tok {
			t.Fatalf("%s: token wrong; expected=%q, got=%q", fs.Position(l.Pos()), test.tok, tok)
		}

		if lit != test.lit {
			t.Fatalf("%s: literal wrong; expected=%q, got=%q", fs.Position(l.Pos()), test.lit, lit)
		}
	}

	err := l.Err()
	if err != nil {
		t.Error(err)
	}
}
