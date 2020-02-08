package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/scanner"

	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/token"
)

type ParserErrors struct {
	errors []ParseError
}

func (p *ParserErrors) AddError(errStr string, lexerError bool) {
	p.errors = append(p.errors, ParseError{err: errors.New(errStr), lexerError: lexerError})
}

func (p ParserErrors) Error() string {
	var (
		header   string
		buf      bytes.Buffer
		writeErr error
	)

	for _, err := range p.errors {
		if err.lexerError {
			header = "lex error: "
		} else {
			header = "parse error: "
		}

		_, writeErr = buf.WriteString(header + err.Error())
		if writeErr != nil {
			panic(writeErr)
		}
	}

	return buf.String()
}

type ParseError struct {
	err        error
	lexerError bool
}

func (p ParseError) Error() string {
	return p.err.Error()
}

type lexer struct {
	err     ParserErrors
	scanner scanner.Scanner

	Statements []ast.Statement
}

func newLexer(in io.Reader) *lexer {
	var s scanner.Scanner

	s.Init(in)
	s.Mode = scanner.ScanInts | scanner.GoTokens

	return &lexer{scanner: s}
}

func (lx *lexer) Lex(yy *yySymType) int {
	tok := lx.scanner.Scan()
	switch tok {
	case scanner.Ident:
		tokText := lx.scanner.TokenText()
		tokType := token.Lookup(tokText)

		switch tokType {
		case token.VAR:
			yy.tok = token.NewToken(token.VAR, lx.scanner.Pos())
			return VAR
		}

		if tokType != token.IDENT {
			panic("invalid token type")
		}

		yy.tok = token.Token{Type: token.IDENT, Pos: lx.scanner.Pos(), Literal: tokText}
		return IDENT
	case scanner.Int:
		yy.tok = token.Token{Type: token.INT, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return INT
	case scanner.Float:
		yy.tok = token.Token{Type: token.FLOAT, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return FLOAT
	case scanner.RawString:
		yy.tok = token.Token{Type: token.RAW_STRING, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return RAW_STRING
	case scanner.EOF:
		return 0
	default:
		switch tok {
		case '=':
			yy.tok = token.Token{Type: token.ASSIGN, Pos: lx.scanner.Pos(), Literal: string(tok)}
			return ASSIGN
		}
	}

	return int(tok)
}

func (lx *lexer) Err() error {
	return lx.err
}

func (lx *lexer) lexerError(s string) {
	lx.err.AddError(lx.scanner.Pos().String()+" "+s, true)
}

func (lx *lexer) lexerErrorf(format string, args ...interface{}) {
	lx.lexerError(fmt.Sprintf(format, args...))
}

func (lx *lexer) Error(s string) {
	lx.err.AddError(lx.scanner.Pos().String()+" "+s, false)
}

func (lx *lexer) Errorf(format string, args ...interface{}) {
	lx.Error(fmt.Sprintf(format, args...))
}
