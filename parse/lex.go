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

	insertSemi bool

	Statements []ast.Statement
}

func newLexer(in io.Reader) *lexer {
	var s scanner.Scanner

	s.Init(in)
	s.Mode = scanner.ScanInts | scanner.GoTokens
	s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

	return &lexer{scanner: s}
}

func (lx *lexer) Lex(yy *yySymType) int {
Scan:
	tok := lx.scanner.Scan()

	if tok == '\n' {
		if lx.insertSemi {
			lx.insertSemi = false
			return SEMICOLON
		}

		goto Scan
	}

	switch tok {
	case scanner.Ident:
		tokText := lx.scanner.TokenText()
		tokType := token.Lookup(tokText)

		switch tokType {
		case token.CONST:
			yy.tok = token.NewToken(token.CONST, lx.scanner.Pos())
			return CONST
		case token.LET:
			yy.tok = token.NewToken(token.LET, lx.scanner.Pos())
			return LET
		case token.VAR:
			yy.tok = token.NewToken(token.VAR, lx.scanner.Pos())
			return VAR
		}

		lx.insertSemi = true
		yy.tok = token.Token{Type: token.IDENT, Pos: lx.scanner.Pos(), Literal: tokText}
		return IDENT
	case scanner.Int:
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.INT, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return INT
	case scanner.Float:
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.FLOAT, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return FLOAT
	case scanner.Char:
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.CHAR, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return CHAR
	case scanner.RawString:
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.RAW_STRING, Pos: lx.scanner.Pos(), Literal: lx.scanner.TokenText()}
		return RAW_STRING
	case '+':
		yy.tok = token.Token{Type: token.ADD, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return ADD
	case '-':
		yy.tok = token.Token{Type: token.SUB, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return SUB
	case '*':
		if lx.scanner.Peek() == '*' {
			yy.tok = token.Token{Type: token.EXP, Pos: lx.scanner.Pos(), Literal: "**"}
			return EXP
		}

		yy.tok = token.Token{Type: token.MUL, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return MUL
	case '/':
		yy.tok = token.Token{Type: token.QUO, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return QUO
	case '%':
		yy.tok = token.Token{Type: token.REM, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return REM

	case '=':
		yy.tok = token.Token{Type: token.ASSIGN, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return ASSIGN

	case '(':
		yy.tok = token.Token{Type: token.LPAREN, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return LPAREN
	case '[':
		yy.tok = token.Token{Type: token.LBRACK, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return LBRACK
	case '{':
		yy.tok = token.Token{Type: token.LBRACE, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return LBRACE
	case ',':
		yy.tok = token.Token{Type: token.COMMA, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return COMMA
	case '.':
		yy.tok = token.Token{Type: token.PERIOD, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return PERIOD
	case ')':
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.RPAREN, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return RPAREN
	case ']':
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.RBRACK, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return RBRACK
	case '}':
		lx.insertSemi = true
		yy.tok = token.Token{Type: token.RBRACE, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return RBRACE
	case ';':
		lx.insertSemi = false
		yy.tok = token.Token{Type: token.SEMICOLON, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return SEMICOLON
	case ':':
		yy.tok = token.Token{Type: token.COLON, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return COLON
	case scanner.EOF:
		if lx.insertSemi {
			lx.insertSemi = false
			return SEMICOLON
		}

		return 0
	}

	return int(tok)
}

func (lx *lexer) Err() error {
	if lx.err.errors != nil {
		return lx.err
	}

	return nil
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
