package parser_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/parser"
	"github.com/capnspacehook/rose/token"

	"github.com/stretchr/testify/require"
)

func TestConstDecls(t *testing.T) {
	expectParse(t, "const i = 5", func(p pfn) []ast.Stmt {
		return stmts(
			constDecl(p(1, 1), 0, 0, valueSpec(
				idents(ident(p(1, 7), "i")), nil, exprs(intLit(p(1, 11), "5")),
			)))
	})

	expectParse(t, "const a, b = 4, 2", func(p pfn) []ast.Stmt {
		return stmts(
			constDecl(p(1, 1), 0, 0,
				valueSpec(
					idents(
						ident(p(1, 7), "a"),
						ident(p(1, 10), "b"),
					),
					nil,
					exprs(
						intLit(p(1, 14), "4"),
						intLit(p(1, 17), "2"),
					),
				)))
	})

	expectParse(t, `const (
	a, b = 4, 2
)`, func(p pfn) []ast.Stmt {
		return stmts(
			constDecl(p(1, 1), p(1, 7), p(3, 1),
				valueSpec(
					idents(
						ident(p(2, 2), "a"),
						ident(p(2, 5), "b"),
					),
					nil,
					exprs(
						intLit(p(2, 9), "4"),
						intLit(p(2, 12), "2"),
					),
				)))
	})

	expectParse(t, `const (
	a = 4
	b = 2
)`, func(p pfn) []ast.Stmt {
		return stmts(
			constDecl(p(1, 1), p(1, 7), p(4, 1),
				valueSpec(
					idents(
						ident(p(2, 2), "a"),
					),
					nil,
					exprs(
						intLit(p(2, 6), "4"),
					)),
				valueSpec(
					idents(
						ident(p(3, 2), "b"),
					),
					nil,
					exprs(
						intLit(p(3, 6), "2"),
					),
				)))
	})

	expectParseError(t, "const a, b = 'v'", "<input>:1:7: missing value in const declaration")
	expectParseError(t, "const a = 'v', 0.9", "<input>:1:7: extra expression in const declaration")
}

func TestVarDecls(t *testing.T) {
	expectParse(t, "var a float", func(p pfn) []ast.Stmt {
		return stmts(
			varDecl(p(1, 1), 0, 0, valueSpec(
				idents(ident(p(1, 5), "a")), ident(p(1, 7), "float"), nil,
			)))
	})

	expectParse(t, "var a, b int", func(p pfn) []ast.Stmt {
		return stmts(
			varDecl(p(1, 1), 0, 0, valueSpec(
				idents(
					ident(p(1, 5), "a"),
					ident(p(1, 8), "b"),
				),
				ident(p(1, 10), "int"),
				nil,
			)))
	})

	expectParse(t, `var (
	a, b string
)`, func(p pfn) []ast.Stmt {
		return stmts(
			varDecl(p(1, 1), p(1, 5), p(3, 1),
				valueSpec(
					idents(
						ident(p(2, 2), "a"),
						ident(p(2, 5), "b"),
					),
					ident(p(2, 7), "string"),
					nil,
				)))
	})

	expectParse(t, `var (
	a float
	b rune
)`, func(p pfn) []ast.Stmt {
		return stmts(
			varDecl(p(1, 1), p(1, 5), p(4, 1),
				valueSpec(
					idents(
						ident(p(2, 2), "a"),
					),
					ident(p(2, 4), "float"),
					nil,
				),
				valueSpec(
					idents(
						ident(p(3, 2), "b"),
					),
					ident(p(3, 4), "rune"),
					nil,
				)))
	})

	// TODO: fix error msg
	expectParseError(t, "var foo, bar int, string", "<input>:1:17: expected ';', found ','")
	expectParseError(t, `var foo = "howdy"`, "<input>:1:5: initialization is not allowed in a var declaration")
}

type pfn func(int, int) token.Pos        // position conversion function
type expectedFn func(pos pfn) []ast.Stmt // callback function to return expected results

func expectParse(t *testing.T, input string, fn expectedFn) {
	fset := token.NewFileSet()
	testFile := fset.AddFile("", -1, len(input))

	var off int
	for {
		o := strings.IndexRune(input[off:], '\n')
		if o == -1 {
			break
		}
		testFile.AddLine(o + off + 1)
		off += o + 1
	}

	actual, err := parser.ParseFile(testFile, strings.NewReader(input))
	require.NoError(t, err)

	expected := fn(func(line, column int) token.Pos {
		return token.Pos(int(testFile.LineStart(line)) + (column - 1))
	})
	require.Equal(t, len(expected), len(actual.Stmts))

	for i := 0; i < len(expected); i++ {
		equalStmt(t, expected[i], actual.Stmts[i])
	}
}

func expectParseError(t *testing.T, input, expectedErr string) {
	fset := token.NewFileSet()
	testFile := fset.AddFile("", -1, len(input))

	_, err := parser.ParseFile(testFile, strings.NewReader(input))
	require.EqualError(t, err, expectedErr)
}

func stmts(s ...ast.Stmt) []ast.Stmt {
	return s
}

func constDecl(pos, lParen, rParen token.Pos, specs ...ast.Spec) ast.Stmt {
	return genDecl(token.CONST, pos, lParen, rParen, specs)
}

func letDecl(pos, lParen, rParen token.Pos, specs ...ast.Spec) ast.Stmt {
	return genDecl(token.LET, pos, lParen, rParen, specs)
}

func varDecl(pos, lParen, rParen token.Pos, specs ...ast.Spec) ast.Stmt {
	return genDecl(token.VAR, pos, lParen, rParen, specs)
}

func genDecl(tok token.Token, pos, lParen, rParen token.Pos, specs []ast.Spec) ast.Stmt {
	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			TokPos: pos,
			Tok:    tok,
			Lparen: lParen,
			Rparen: rParen,
			Specs:  specs,
		},
	}
}

func valueSpec(names []*ast.Ident, typ ast.Expr, values []ast.Expr) *ast.ValueSpec {
	return &ast.ValueSpec{
		Names:  names,
		Type:   typ,
		Values: values,
	}
}

func idents(list ...*ast.Ident) []*ast.Ident {
	return list
}

func ident(pos token.Pos, name string) *ast.Ident {
	return &ast.Ident{
		NamePos: pos,
		Name:    name,
	}
}

func exprs(list ...ast.Expr) []ast.Expr {
	return list
}

func intLit(pos token.Pos, val string) *ast.BasicLit {
	return basicLit(pos, token.INT, val)
}

func basicLit(pos token.Pos, kind token.Token, val string) *ast.BasicLit {
	return &ast.BasicLit{
		ValuePos: pos,
		Kind:     kind,
		Value:    val,
	}
}

func equalStmt(t *testing.T, expected, actual ast.Stmt) {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		require.Nil(t, actual, "expected nil, but got not nil")
		return
	}
	require.NotNil(t, actual, "expected not nil, but got nil")
	require.IsType(t, expected, actual)

	switch expected := expected.(type) {
	case *ast.DeclStmt:
		equalDecl(t, expected.Decl, actual.(*ast.DeclStmt).Decl)
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}

func equalDecl(t *testing.T, expected, actual ast.Decl) {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		require.Nil(t, actual, "expected nil, but got not nil")
		return
	}
	require.NotNil(t, actual, "expected not nil, but got nil")
	require.IsType(t, expected, actual)

	switch expected := expected.(type) {
	case *ast.GenDecl:
		require.Equal(t, expected.TokPos, actual.(*ast.GenDecl).TokPos)
		require.Equal(t, expected.Tok, actual.(*ast.GenDecl).Tok)
		require.Equal(t, expected.Lparen, actual.(*ast.GenDecl).Lparen)
		require.Equal(t, expected.Rparen, actual.(*ast.GenDecl).Rparen)

		require.Equal(t, len(expected.Specs), len(actual.(*ast.GenDecl).Specs))
		for i := 0; i < len(expected.Specs); i++ {
			equalSpec(t, expected.Specs[i], actual.(*ast.GenDecl).Specs[i])
		}
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}

func equalSpec(t *testing.T, expected, actual ast.Spec) {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		require.Nil(t, actual, "expected nil, but got not nil")
		return
	}
	require.NotNil(t, actual, "expected not nil, but got nil")
	require.IsType(t, expected, actual)

	switch expected := expected.(type) {
	case *ast.ValueSpec:
		require.Equal(t, len(expected.Names), len(actual.(*ast.ValueSpec).Names))
		for i := 0; i < len(expected.Names); i++ {
			equalIdent(t, expected.Names[i], actual.(*ast.ValueSpec).Names[i])
		}

		equalExpr(t, expected.Type, actual.(*ast.ValueSpec).Type)

		require.Equal(t, len(expected.Values), len(actual.(*ast.ValueSpec).Values))
		for i := 0; i < len(expected.Values); i++ {
			equalExpr(t, expected.Values[i], actual.(*ast.ValueSpec).Values[i])
		}
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}

func equalIdent(t *testing.T, expected, actual *ast.Ident) {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		require.Nil(t, actual, "expected nil, but got not nil")
		return
	}
	require.NotNil(t, actual, "expected not nil, but got nil")

	require.Equal(t, expected.NamePos, actual.NamePos)
	require.Equal(t, expected.Name, actual.Name)
}

func equalExpr(t *testing.T, expected, actual ast.Expr) {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		require.Nil(t, actual, "expected nil, but got not nil")
		return
	}
	require.NotNil(t, actual, "expected not nil, but got nil")
	require.IsType(t, expected, actual)

	switch expected := expected.(type) {
	case *ast.Ident:
		equalIdent(t, expected, actual.(*ast.Ident))
	case *ast.BasicLit:
		require.Equal(t, expected.ValuePos, actual.(*ast.BasicLit).ValuePos)
		require.Equal(t, expected.Kind, actual.(*ast.BasicLit).Kind)
		require.Equal(t, expected.Value, actual.(*ast.BasicLit).Value)
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}
