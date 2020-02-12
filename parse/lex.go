package parse

import (
	"fmt"
	"io"
	"text/scanner"

	"github.com/capnspacehook/rose/token"
)

type lexer struct {
	file       *token.File // source file handle
	errors     ErrorList   // lexing errors
	insertSemi bool        // insert a semicolon before next newline

	scanner scanner.Scanner // scanner that does much of the heavy lifting
}

func (lx *lexer) Init(file *token.File, src io.Reader, srcLen int) {
	if file.Size() != srcLen {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), srcLen))
	}

	// explicitly initialize all fields since a scanner may be reused
	lx.file = file
	lx.errors = ErrorList{}
	lx.insertSemi = false

	lx.scanner.Init(src)
	lx.scanner.Mode = scanner.ScanInts | scanner.GoTokens
	lx.scanner.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	lx.scanner.Error = func(s *scanner.Scanner, msg string) {
		lx.error(msg)
	}
}

func (lx *lexer) Lex() (pos token.Pos, tok token.Token, lit string) {
lexAgain:
	ch := lx.scanner.Scan()
	pos = lx.currentPos()

	if ch == '\n' {
		if lx.insertSemi {
			lx.insertSemi = false
			return pos, token.SEMICOLON, "\n"
		}

		goto lexAgain
	}

	switch ch {
	case scanner.Ident:
		tok = token.IDENT
		lit = lx.scanner.TokenText()
	case scanner.Int:
		lx.insertSemi = true
		tok = token.INT
		lit = lx.scanner.TokenText()
	case scanner.Float:
		lx.insertSemi = true
		tok = token.FLOAT
		lit = lx.scanner.TokenText()
	case scanner.Char:
		lx.insertSemi = true
		tok = token.CHAR
		lit = lx.scanner.TokenText()
	case scanner.String: // TODO: handle string expressions
		lx.insertSemi = true
		tok = token.STRING
		lit = lx.scanner.TokenText()
	case scanner.RawString:
		lx.insertSemi = true
		tok = token.RAW_STRING
		lit = lx.scanner.TokenText()
	case '+':
		pCh := lx.scanner.Peek()
		if pCh == '+' {
			lx.scanner.Next()
			return lx.currentPos(), token.INC, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.ADD_ASSIGN, ""
		}

		tok = token.ADD
	case '-':
		pCh := lx.scanner.Peek()
		if pCh == '-' {
			lx.scanner.Next()
			return lx.currentPos(), token.DEC, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.SUB_ASSIGN, ""
		}

		tok = token.SUB
	case '*':
		pCh := lx.scanner.Peek()
		if pCh == '*' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				return lx.currentPos(), token.EXP_ASSIGN, ""
			}

			return lx.currentPos(), token.EXP, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.MUL_ASSIGN, ""
		}

		tok = token.MUL
	case '/':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.QUO_ASSIGN, ""
		}

		tok = token.QUO
	case '%':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.REM_ASSIGN, ""
		}

		tok = token.REM
	case '&':
		pCh := lx.scanner.Peek()
		if pCh == '^' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				return lx.currentPos(), token.AND_NOT_ASSIGN, ""
			}

			return lx.currentPos(), token.AND_NOT, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.AND_ASSIGN, ""
		}

		tok = token.AND
	case '|':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.OR_ASSIGN, ""
		}

		tok = token.OR
	case '^':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.XOR_ASSIGN, ""
		}

		tok = token.XOR
	case '~':
		tok = token.INVT
	case '<':
		pCh := lx.scanner.Peek()
		if pCh == '<' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				return lx.currentPos(), token.SHL_ASSIGN, ""
			}

			return lx.currentPos(), token.SHL, ""
		} else if pCh == '-' {
			lx.scanner.Next()
			return lx.currentPos(), token.ARROW, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.LEQ, ""
		}

		tok = token.LSS
	case '>':
		pCh := lx.scanner.Peek()
		if pCh == '>' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '=' {
				lx.scanner.Next()
				return lx.currentPos(), token.SHR_ASSIGN, ""
			}

			return lx.currentPos(), token.SHR, ""
		} else if pCh == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.GEQ, ""
		}

		tok = token.GTR
	case '=':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.EQL, ""
		}

		tok = token.ASSIGN
	case '(':
		tok = token.LPAREN
	case '[':
		tok = token.LBRACK
	case '{':
		tok = token.LBRACE
	case ',':
		tok = token.COMMA
	case '.':
		if lx.scanner.Peek() == '.' {
			lx.scanner.Next()
			if lx.scanner.Peek() == '.' {
				lx.scanner.Next()
				return lx.currentPos(), token.ELLIPSIS, ""
			}
		}

		tok = token.PERIOD
	case ')':
		lx.insertSemi = true
		tok = token.RPAREN
	case ']':
		lx.insertSemi = true
		tok = token.RBRACK
	case '}':
		lx.insertSemi = true
		tok = token.RBRACE
	case ';':
		lx.insertSemi = false
		tok = token.SEMICOLON
	case ':':
		tok = token.COLON
	case '?':
		tok = token.QUES
	case '!':
		if lx.scanner.Peek() == '=' {
			lx.scanner.Next()
			return lx.currentPos(), token.NEQ, ""
		}

		tok = token.EXCLM
	case scanner.EOF:
		if lx.insertSemi {
			lx.insertSemi = false
			tok = token.SEMICOLON
		} else {
			tok = token.EOF
		}
	default:
		lx.errorf("illegal character %#U", ch)
	}

	return
}

func (lx *lexer) currentPos() token.Pos {
	return lx.file.Pos(lx.scanner.Offset)
}

func (lx *lexer) Err() error {
	return lx.errors.Err()
}

func (lx *lexer) error(msg string) {
	lx.errors.Add(lx.scanner.Pos(), msg)
}

func (lx *lexer) errorf(format string, args ...interface{}) {
	lx.error(fmt.Sprintf(format, args...))
}
