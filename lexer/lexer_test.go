package lexer

import (
	"strings"
	"testing"
	"text/scanner"

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
{1: "one", 2: "two"};
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

		{token.IF, "if"},
		{token.INT, "5"},
		{token.LSS, ""},
		{token.INT, "9"},
		{token.LBRACE, ""},
		{token.RETURN, "return"},
		{token.IDENT, "true"},
		{token.SEMI, "\n"},
		{token.RBRACE, ""},
		{token.ELSE, "else"},
		{token.LBRACE, ""},
		{token.RETURN, "return"},
		{token.IDENT, "false"},
		{token.SEMI, "\n"},
		{token.RBRACE, ""},
		{token.SEMI, "\n"},

		{token.LBRACK, ""},
		{token.INT, "1"},
		{token.COMMA, ""},
		{token.INT, "2"},
		{token.COMMA, ""},
		{token.INT, "3"},
		{token.RBRACK, ""},
		{token.SEMI, "\n"},
		{token.LBRACE, ""},
		{token.INT, "1"},
		{token.COLON, ""},
		{token.STRING, `"one"`},
		{token.COMMA, ""},
		{token.INT, "2"},
		{token.COLON, ""},
		{token.STRING, `"two"`},
		{token.RBRACE, ""},
		{token.SEMI, ";"},
		{token.IDENT, "l"},
		{token.LBRACK, ""},
		{token.INT, "2"},
		{token.COLON, ""},
		{token.INT, "5"},
		{token.RBRACK, ""},
		{token.SEMI, "\n"},
		{token.EOF, ""},
	}

	var l Lexer
	var errs ErrorList
	fs := token.NewFileSet()
	eh := func(pos token.Position, msg string) { errs.Add(scanner.Position(pos), msg) }

	l.Init(fs.AddFile("", -1, len(input)), strings.NewReader(input), eh, true)
	for _, test := range tests {
		pos, tok, lit := l.Lex()
		if tok != test.tok {
			t.Fatalf("%s: token wrong; expected=%q, got=%q", fs.Position(pos), test.tok, tok)
		}

		if lit != test.lit {
			t.Fatalf("%s: literal wrong; expected=%q, got=%q", fs.Position(pos), test.lit, lit)
		}
	}

	if err := errs.Err(); err != nil {
		t.Error(err)
	}
}
