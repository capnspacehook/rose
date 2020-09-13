// Copyright 2009 The Go Authors. All rights reserved.

// Package token defines constants representing the lexical tokens of the Rose
// programming language and basic operations on tokens (printing, predicates).

package token

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

var (
	operatorWords = make(map[string]Token)
	keywords      = make(map[string]Token)
)

func init() {
	operatorWords[tokens[LAND]] = LAND
	operatorWords[tokens[LOR]] = LOR
	operatorWords[tokens[NOT]] = NOT

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
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
	INVT    // ~
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=
	EXP_ASSIGN // **=

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

	RPAREN // )
	RBRACK // ]
	RBRACE // }
	SEMI   // ;
	COLON  // :

	QUES  // ?
	EXCLM // !
	operator_end

	keyword_beg
	// Keywords
	BREAK
	CONST
	CONTINUE
	FALLTHROUGH
	FN
	LET
	RETURN
	VAR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:      "IDENT",
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
	INVT:    "~",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",
	EXP_ASSIGN: "**",

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

	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",
	SEMI:   ";",
	COLON:  ":",

	QUES:  "?",
	EXCLM: "!",

	BREAK:       "break",
	CONST:       "const",
	CONTINUE:    "continue",
	FALLTHROUGH: "fallthrough",
	FN:          "fn",
	LET:         "let",
	RETURN:      "return",
	VAR:         "var",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 7
	HighestPrec = 8
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
// TODO: add QUES and EXCLM
func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	case EXP:
		return 6
	}

	return LowestPrec
}

// Lookup maps an identifier to its keyword token, operator token
// (if LAND, LOR, or NOT), or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	if tok, is_operator := operatorWords[ident]; is_operator {
		return tok
	}

	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsExported reports whether name starts with an upper-case letter.
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsKeyword reports whether name is a Go keyword, such as "func" or "return".
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}

	return name != "" && !IsKeyword(name)
}
