package token

import (
	"text/scanner"
)

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	IDENT      // main
	INT        // 12345
	FLOAT      // 123.45
	IMAG       // 123.45i
	CHAR       // 'a'
	STRING     // "abc"
	RAW_STRING // `abc`
	literal_end

	operator_beg
	// Operators and punctuation
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	EXP // **

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // and
	LOR   // or
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // not

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	QUES  // ?
	EXCLM // !
	operator_end

	keyword_beg
	// Keywords
	CONST
	LET
	VAR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	INT:        "INT",
	FLOAT:      "FLOAT",
	CHAR:       "CHAR",
	STRING:     "STRING",
	RAW_STRING: "RAW_STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",
	EXP: "**",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "and",
	LOR:   "or",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "not",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	QUES:  "?",
	EXCLM: "!",

	CONST: "const",
	LET:   "let",
	VAR:   "var",
}

type Token struct {
	Type    TokenType
	Pos     scanner.Position
	Literal string
}

func NewToken(tokType TokenType, pos scanner.Position) Token {
	return Token{Type: tokType, Pos: pos, Literal: tokens[tokType]}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) TokenType {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	return IDENT
}
