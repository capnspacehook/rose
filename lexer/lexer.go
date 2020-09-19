package lexer

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/capnspacehook/rose/token"
)

type Lexer struct {
	file       *token.File  // source file handle
	errh       ErrorHandler // error reporting
	insertSemi bool         // insert a semicolon before next newline

	strBuf  strings.Builder
	scanner scanner.Scanner // scanner that does much of the heavy lifting
}

// An ErrorHandler may be provided to Lexer.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// position and an error message. The position points to the beginning of
// the offending token.
type ErrorHandler func(pos token.Position, msg string)

func (lx *Lexer) Init(file *token.File, src io.Reader, errh ErrorHandler, scanComments bool) {
	// explicitly initialize all fields since a scanner may be reused
	lx.file = file
	lx.insertSemi = false

	lx.scanner.Init(src)
	lx.scanner.Filename = file.Name()
	lx.scanner.Mode = scanner.ScanInts | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings
	lx.scanner.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	if errh != nil {
		lx.scanner.Error = func(s *scanner.Scanner, msg string) {
			errh(token.Position(s.Pos()), msg)
		}
	}
}

func (lx *Lexer) Lex() (pos token.Pos, tok token.Token, lit string) {
lexAgain:
	ch := lx.scanner.Scan()
	pos = lx.currentPos()

	insertSemi := false
	switch ch {
	case scanner.Int:
		insertSemi = true
		tok = token.INT
		lit = lx.scanner.TokenText()
	case scanner.Float:
		insertSemi = true
		tok = token.FLOAT
		lit = lx.scanner.TokenText()
	case scanner.Char:
		insertSemi = true
		tok = token.CHAR
		lit = lx.scanner.TokenText()
	case scanner.String: // TODO: (capnspacehook) handle string expressions
		insertSemi = true
		tok = token.STRING
		lit = lx.scanner.TokenText()
	case scanner.RawString:
		insertSemi = true
		tok = token.RAW_STRING
		lit = lx.scanner.TokenText()
	case '+':
		switch lx.scanner.Peek() {
		case '+':
			insertSemi = true
			lx.scanner.Next()
			tok = token.INC
		case '=':
			lx.scanner.Next()
			tok = token.ADD_ASSIGN
		default:
			tok = token.ADD
		}
	case '-':
		switch lx.scanner.Peek() {
		case '-':
			insertSemi = true
			lx.scanner.Next()
			tok = token.DEC
		case '=':
			lx.scanner.Next()
			tok = token.SUB_ASSIGN
		default:
			tok = token.SUB
		}
	case '*':
		switch lx.scanner.Peek() {
		case '*':
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				tok = token.EXP_ASSIGN
			} else {
				tok = token.EXP
			}
		case '=':
			lx.scanner.Next()
			tok = token.MUL_ASSIGN
		default:
			tok = token.MUL
		}
	case '/':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.QUO_ASSIGN
		default:
			tok = token.QUO
		}
	case '%':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.REM_ASSIGN
		default:
			tok = token.REM
		}
	case '&':
		switch lx.scanner.Peek() {
		case '^':
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				tok = token.AND_NOT_ASSIGN
			} else {
				tok = token.AND_NOT
			}
		case '=':
			lx.scanner.Next()
			tok = token.AND_ASSIGN
		default:
			tok = token.AND
		}
	case '|':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.OR_ASSIGN
		default:
			tok = token.OR
		}
	case '^':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.XOR_ASSIGN
		default:
			tok = token.XOR
		}
	case '~':
		tok = token.INVT
	case '<':
		switch lx.scanner.Peek() {
		case '<':
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				tok = token.SHL_ASSIGN
			} else {
				tok = token.SHL
			}
		case '-':
			lx.scanner.Next()
			tok = token.ARROW
		case '=':
			lx.scanner.Next()
			tok = token.LEQ
		default:
			tok = token.LSS
		}
	case '>':
		switch lx.scanner.Peek() {
		case '>':
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				tok = token.SHR_ASSIGN
			} else {
				tok = token.SHR
			}
		case '=':
			lx.scanner.Next()
			tok = token.GEQ
		default:
			tok = token.GTR
		}
	case '=':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.EQL
		default:
			tok = token.ASSIGN
		}
	case '(':
		tok = token.LPAREN
	case '[':
		tok = token.LBRACK
	case '{':
		tok = token.LBRACE
	case ',':
		tok = token.COMMA
	case '.':
		switch lx.scanner.Peek() {
		case '.':
			lx.scanner.Next()
			if lx.scanner.Peek() == '.' {
				lx.scanner.Next()
				tok = token.ELLIPSIS
			}
		default:
			tok = token.PERIOD
		}
	case ')':
		insertSemi = true
		tok = token.RPAREN
	case ']':
		insertSemi = true
		tok = token.RBRACK
	case '}':
		insertSemi = true
		tok = token.RBRACE
	case ';':
		insertSemi = false
		tok = token.SEMI
		lit = ";"
	case ':':
		tok = token.COLON
	case '?':
		tok = token.QUES
	case '!':
		switch lx.scanner.Peek() {
		case '=':
			lx.scanner.Next()
			tok = token.NEQ
		default:
			tok = token.EXCLM
		}
	case '\n':
		if lx.insertSemi {
			lx.insertSemi = false
			return pos, token.SEMI, "\n"
		}

		goto lexAgain
	case scanner.EOF:
		if lx.insertSemi {
			lx.insertSemi = false
			return pos, token.SEMI, ""
		}

		tok = token.EOF
	default:
		if lx.isIdentRune(ch, true) {
			insertSemi = true
			lit = lx.scanIdentifier(ch)
			if len(lit) > 1 {
				// keywords are longer than one letter - avoid lookup otherwise
				switch tok = token.Lookup(lit); tok {
				case token.IDENT, token.BREAK, token.CONTINUE, token.FALLTHROUGH, token.RETURN:
					insertSemi = true
				}
			} else {
				insertSemi = true
				tok = token.IDENT
			}
		} else {
			lx.scanner.Error(&lx.scanner, fmt.Sprintf("illegal character %#U", ch))
		}
	}

	lx.insertSemi = insertSemi

	return
}

func (lx *Lexer) currentPos() token.Pos {
	return lx.file.Pos(lx.scanner.Offset)
}

func (lx *Lexer) scanIdentifier(ch rune) string {
	// we know the zero'th rune is OK; start scanning at the next one
	lx.strBuf.WriteRune(ch)
	for lx.isIdentRune(lx.scanner.Peek(), false) {
		lx.strBuf.WriteRune(lx.scanner.Scan())
	}
	defer lx.strBuf.Reset()

	return lx.strBuf.String()
}

func (lx *Lexer) isIdentRune(ch rune, firstRune bool) bool {
	return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && !firstRune
}
