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

func (p *ParserErrors) AddError(errStr string, pos scanner.Position, lexerError bool) {
	p.errors = append(p.errors, ParseError{err: errors.New(errStr), pos: pos, lexerError: lexerError})
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

		_, writeErr = buf.WriteString(header + err.pos.String() + ": " + err.Error() + "\n")
		if writeErr != nil {
			panic(writeErr)
		}
	}

	return buf.String()
}

type ParseError struct {
	err        error
	pos        scanner.Position
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
	var l lexer

	l.scanner.Init(in)
	l.scanner.Mode = scanner.ScanInts | scanner.GoTokens
	l.scanner.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	l.scanner.Error = func(s *scanner.Scanner, msg string) {
		l.Error(msg)
	}

	return &l
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
		pTok := lx.scanner.Peek()
		if pTok == '+' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.INC, Pos: lx.scanner.Pos(), Literal: "++"}
			return INC
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.ADD_ASSIGN, Pos: lx.scanner.Pos(), Literal: "+="}
			return ADD_ASSIGN
		}

		yy.tok = token.Token{Type: token.ADD, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return ADD
	case '-':
		pTok := lx.scanner.Peek()
		if pTok == '-' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.DEC, Pos: lx.scanner.Pos(), Literal: "--"}
			return DEC
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.SUB_ASSIGN, Pos: lx.scanner.Pos(), Literal: "-="}
			return SUB_ASSIGN
		}

		yy.tok = token.Token{Type: token.SUB, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return SUB
	case '*':
		pTok := lx.scanner.Peek()
		if pTok == '*' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.EXP, Pos: lx.scanner.Pos(), Literal: "**"}
			return EXP
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.MUL_ASSIGN, Pos: lx.scanner.Pos(), Literal: "*="}
			return MUL_ASSIGN
		}

		yy.tok = token.Token{Type: token.MUL, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return MUL
	case '/':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.QUO_ASSIGN, Pos: lx.scanner.Pos(), Literal: "/="}
			return QUO_ASSIGN
		}

		yy.tok = token.Token{Type: token.QUO, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return QUO
	case '%':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.REM_ASSIGN, Pos: lx.scanner.Pos(), Literal: "%="}
			return REM_ASSIGN
		}

		yy.tok = token.Token{Type: token.REM, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return REM
	case '&':
		pTok := lx.scanner.Peek()
		if pTok == '^' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				yy.tok = token.Token{Type: token.AND_NOT_ASSIGN, Pos: lx.scanner.Pos(), Literal: "&^"}
				return AND_NOT_ASSIGN
			}

			yy.tok = token.Token{Type: token.AND_NOT, Pos: lx.scanner.Pos(), Literal: "&^"}
			return AND_NOT
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.AND_ASSIGN, Pos: lx.scanner.Pos(), Literal: "&="}
			return AND_ASSIGN
		}

		yy.tok = token.Token{Type: token.AND, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return AND
	case '|':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.OR_ASSIGN, Pos: lx.scanner.Pos(), Literal: "|="}
			return OR_ASSIGN
		}

		yy.tok = token.Token{Type: token.OR, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return OR
	case '^':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.XOR_ASSIGN, Pos: lx.scanner.Pos(), Literal: "^="}
			return XOR_ASSIGN
		}

		yy.tok = token.Token{Type: token.XOR, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return XOR
	case '<':
		pTok := lx.scanner.Peek()
		if pTok == '<' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				yy.tok = token.Token{Type: token.SHL_ASSIGN, Pos: lx.scanner.Pos(), Literal: "<<="}
				return SHL_ASSIGN
			}

			yy.tok = token.Token{Type: token.SHL, Pos: lx.scanner.Pos(), Literal: "<<"}
			return SHL
		} else if pTok == '-' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.ARROW, Pos: lx.scanner.Pos(), Literal: "<-"}
			return ARROW
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.LEQ, Pos: lx.scanner.Pos(), Literal: "<="}
			return LEQ
		}

		yy.tok = token.Token{Type: token.LSS, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return LSS
	case '>':
		pTok := lx.scanner.Peek()
		if pTok == '>' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				yy.tok = token.Token{Type: token.SHR_ASSIGN, Pos: lx.scanner.Pos(), Literal: ">>="}
				return SHR_ASSIGN
			}
			yy.tok = token.Token{Type: token.SHR, Pos: lx.scanner.Pos(), Literal: ">>"}
			return SHR
		} else if pTok == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.GEQ, Pos: lx.scanner.Pos(), Literal: ">="}
			return GEQ
		}

		yy.tok = token.Token{Type: token.GTR, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return GTR
	case '=':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.EQL, Pos: lx.scanner.Pos(), Literal: "=="}
			return EQL
		}

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
		if lx.scanner.Peek() == '.' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '.' {
				lx.scanner.Next()
				yy.tok = token.Token{Type: token.ELLIPSIS, Pos: lx.scanner.Pos(), Literal: "..."}
				return ELLIPSIS
			}
		}

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
	case '?':
		yy.tok = token.Token{Type: token.QUES, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return QUES
	case '!':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			yy.tok = token.Token{Type: token.NEQ, Pos: lx.scanner.Pos(), Literal: "!="}
			return NEQ
		}

		yy.tok = token.Token{Type: token.EXCLM, Pos: lx.scanner.Pos(), Literal: string(tok)}
		return EXCLM
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

func (lx *lexer) lexerError(msg string) {
	lx.err.AddError(msg, lx.scanner.Pos(), true)
}

func (lx *lexer) lexerErrorf(format string, args ...interface{}) {
	lx.lexerError(fmt.Sprintf(format, args...))
}

func (lx *lexer) Error(msg string) {
	lx.err.AddError(msg, lx.scanner.Pos(), false)
}

func (lx *lexer) Errorf(format string, args ...interface{}) {
	lx.Error(fmt.Sprintf(format, args...))
}
