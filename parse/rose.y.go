// Code generated by goyacc -o rose.y.go rose.y. DO NOT EDIT.

//line rose.y:2
package parse

import __yyfmt__ "fmt"

//line rose.y:2

import (
	"fmt"
	"strconv"

	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/token"
)

//line rose.y:13
type yySymType struct {
	yys      int
	stmt     ast.Statement
	stmtlist []ast.Statement
	expr     ast.Expression
	exprlist []ast.Expression
	typename *ast.TypeName
	tok      token.Token
}

const IDENT = 57346
const INT = 57347
const FLOAT = 57348
const CHAR = 57349
const STRING = 57350
const RAW_STRING = 57351
const ADD = 57352
const SUB = 57353
const MUL = 57354
const QUO = 57355
const REM = 57356
const EXP = 57357
const AND = 57358
const OR = 57359
const XOR = 57360
const SHL = 57361
const SHR = 57362
const AND_NOT = 57363
const ADD_ASSIGN = 57364
const SUB_ASSIGN = 57365
const MUL_ASSIGN = 57366
const QUO_ASSIGN = 57367
const REM_ASSIGN = 57368
const EXP_ASSIGN = 57369
const AND_ASSIGN = 57370
const OR_ASSIGN = 57371
const XOR_ASSIGN = 57372
const SHL_ASSIGN = 57373
const SHR_ASSIGN = 57374
const AND_NOT_ASSIGN = 57375
const LAND = 57376
const LOR = 57377
const ARROW = 57378
const INC = 57379
const DEC = 57380
const EQL = 57381
const LSS = 57382
const GTR = 57383
const ASSIGN = 57384
const NOT = 57385
const NEQ = 57386
const LEQ = 57387
const GEQ = 57388
const ELLIPSIS = 57389
const LPAREN = 57390
const LBRACK = 57391
const LBRACE = 57392
const COMMA = 57393
const PERIOD = 57394
const RPAREN = 57395
const RBRACK = 57396
const RBRACE = 57397
const SEMICOLON = 57398
const COLON = 57399
const QUES = 57400
const EXCLM = 57401
const CONST = 57402
const LET = 57403
const VAR = 57404

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"INT",
	"FLOAT",
	"CHAR",
	"STRING",
	"RAW_STRING",
	"ADD",
	"SUB",
	"MUL",
	"QUO",
	"REM",
	"EXP",
	"AND",
	"OR",
	"XOR",
	"SHL",
	"SHR",
	"AND_NOT",
	"ADD_ASSIGN",
	"SUB_ASSIGN",
	"MUL_ASSIGN",
	"QUO_ASSIGN",
	"REM_ASSIGN",
	"EXP_ASSIGN",
	"AND_ASSIGN",
	"OR_ASSIGN",
	"XOR_ASSIGN",
	"SHL_ASSIGN",
	"SHR_ASSIGN",
	"AND_NOT_ASSIGN",
	"LAND",
	"LOR",
	"ARROW",
	"INC",
	"DEC",
	"EQL",
	"LSS",
	"GTR",
	"ASSIGN",
	"NOT",
	"NEQ",
	"LEQ",
	"GEQ",
	"ELLIPSIS",
	"LPAREN",
	"LBRACK",
	"LBRACE",
	"COMMA",
	"PERIOD",
	"RPAREN",
	"RBRACK",
	"RBRACE",
	"SEMICOLON",
	"COLON",
	"QUES",
	"EXCLM",
	"CONST",
	"LET",
	"VAR",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 143

var yyAct = [...]int{

	28, 11, 77, 78, 81, 82, 83, 84, 87, 79,
	80, 85, 86, 88, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 66, 67, 7, 60,
	60, 71, 73, 75, 60, 96, 72, 74, 76, 94,
	92, 90, 29, 27, 26, 98, 77, 78, 81, 82,
	83, 84, 87, 79, 80, 85, 86, 88, 25, 91,
	59, 1, 93, 31, 95, 70, 97, 63, 61, 69,
	66, 67, 58, 68, 89, 71, 73, 75, 65, 39,
	72, 74, 76, 32, 9, 10, 8, 62, 64, 30,
	6, 5, 4, 99, 3, 100, 2, 101, 38, 41,
	42, 43, 44, 45, 33, 34, 0, 0, 0, 0,
	0, 0, 36, 13, 14, 15, 16, 17, 18, 19,
	20, 21, 22, 23, 24, 0, 0, 0, 0, 0,
	37, 0, 0, 12, 0, 0, 0, 35, 0, 0,
	0, 0, 40,
}
var yyPact = [...]int{

	-1000, -1000, 24, -55, -1000, -1000, -1000, 91, 54, 40,
	39, -1000, 94, 94, 94, 94, 94, 94, 94, 94,
	94, 94, 94, 94, 94, 30, 26, 25, 36, -1000,
	-1000, 94, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	94, -1000, -1000, -1000, -1000, -1000, 36, 36, 36, 36,
	36, 36, 36, 36, 36, 36, 36, 36, 94, -2,
	-1000, 94, -3, 94, -7, 94, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-8, 36, 94, 36, 94, 36, 94, 36, -1000, 36,
	36, 36,
}
var yyPgo = [...]int{

	0, 96, 94, 92, 91, 90, 60, 0, 42, 89,
	83, 79, 78, 73, 69, 65, 63, 61,
}
var yyR1 = [...]int{

	0, 17, 1, 1, 2, 2, 2, 4, 4, 6,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 5, 5, 5, 5, 7, 7, 8,
	8, 9, 10, 10, 10, 11, 11, 11, 11, 11,
	12, 12, 12, 12, 12, 13, 13, 13, 13, 13,
	13, 14, 14, 14, 14, 15, 15, 15, 15, 15,
	15, 15, 15, 16, 16, 16, 16, 16,
}
var yyR2 = [...]int{

	0, 1, 0, 3, 1, 1, 1, 4, 5, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 4, 5, 4, 5, 1, 3, 1,
	2, 1, 1, 1, 3, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -17, -1, -2, -3, -4, -5, 4, 62, 60,
	61, 56, 42, 22, 23, 24, 25, 26, 27, 28,
	29, 30, 31, 32, 33, 4, 4, 4, -7, -8,
	-9, -16, -10, 10, 11, 43, 18, 36, 4, -11,
	48, 5, 6, 7, 8, 9, -7, -7, -7, -7,
	-7, -7, -7, -7, -7, -7, -7, -7, 42, -6,
	4, 42, -6, 42, -6, -12, 34, 35, -13, -14,
	-15, 39, 44, 40, 45, 41, 46, 10, 11, 17,
	18, 12, 13, 14, 15, 19, 20, 16, 21, -8,
	-7, -7, 42, -7, 42, -7, 42, -7, 53, -7,
	-7, -7,
}
var yyDef = [...]int{

	2, -2, 1, 0, 4, 5, 6, 0, 0, 0,
	0, 3, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 10, 27,
	29, 0, 31, 63, 64, 65, 66, 67, 32, 33,
	0, 35, 36, 37, 38, 39, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22, 0, 0,
	9, 0, 0, 0, 0, 0, 40, 41, 42, 43,
	44, 45, 46, 47, 48, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 58, 59, 60, 61, 62, 30,
	0, 7, 0, 23, 0, 25, 0, 28, 34, 8,
	24, 26,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:56
		{
			yylex.(*lexer).Statements = yyDollar[1].stmtlist
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
//line rose.y:62
		{
			yyVAL.stmtlist = nil
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:64
		{
			yyVAL.stmtlist = append(yyDollar[1].stmtlist, yyDollar[2].stmt)
		}
	case 7:
		yyDollar = yyS[yypt-4 : yypt+1]
//line rose.y:77
		{
			yyVAL.stmt = &ast.VarDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Value: yyDollar[4].expr,
			}
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
//line rose.y:85
		{
			yyVAL.stmt = &ast.VarDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Type:  yyDollar[3].typename,
				Value: yyDollar[5].expr,
			}
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:97
		{
			if _, ok := typeNames[yyDollar[1].tok.Literal]; ok {
				yyVAL.typename = &ast.TypeName{
					Token: yyDollar[1].tok,
				}
			} else {
				yylex.Error(fmt.Sprintf("%q is not a valid type", yyDollar[1].tok.Literal))
			}
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:110
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:118
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:126
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:134
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:142
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:150
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:158
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:166
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:174
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:182
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:190
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:198
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:206
		{
			yyVAL.stmt = &ast.AssignmentStatement{
				Token: yyDollar[2].tok,
				Name:  &ast.Identifier{Token: yyDollar[1].tok},
				Value: yyDollar[3].expr,
			}
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
//line rose.y:217
		{
			yyVAL.stmt = &ast.ConstDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Value: yyDollar[4].expr,
			}
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
//line rose.y:225
		{
			yyVAL.stmt = &ast.ConstDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Type:  yyDollar[3].typename,
				Value: yyDollar[5].expr,
			}
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
//line rose.y:234
		{
			yyVAL.stmt = &ast.ConstDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Value: yyDollar[4].expr,
			}
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
//line rose.y:242
		{
			yyVAL.stmt = &ast.ConstDeclStatement{
				Token: yyDollar[1].tok,
				Name:  &ast.Identifier{Token: yyDollar[2].tok},
				Type:  yyDollar[3].typename,
				Value: yyDollar[5].expr,
			}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:255
		{
			yyVAL.expr = &ast.BinaryExpression{
				Lhs:   yyDollar[1].expr,
				Token: yyDollar[2].tok,
				Rhs:   yyDollar[3].expr,
			}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
//line rose.y:267
		{
			yyVAL.expr = &ast.UnaryExpression{
				Token: yyDollar[1].tok,
				Value: yyDollar[2].expr,
			}
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:281
		{
			if v, ok := boolConsts[yyDollar[1].tok.Literal]; ok {
				yyVAL.expr = &ast.Boolean{
					Token: yyDollar[1].tok,
					Value: v,
				}
			} else if yyDollar[1].tok.Literal == "nil" {
				yyVAL.expr = &ast.Nil{
					Token: yyDollar[1].tok,
				}
			} else {
				yyVAL.expr = &ast.Identifier{
					Token: yyDollar[1].tok,
				}
			}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
//line rose.y:299
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:304
		{
			i, err := strconv.ParseInt(yyDollar[1].tok.Literal, 0, 64)
			if err != nil {
				yylex.(*lexer).Error(err.Error())
			}

			yyVAL.expr = &ast.IntegerLiteral{
				Token: yyDollar[1].tok,
				Value: i,
			}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:316
		{
			f, err := strconv.ParseFloat(yyDollar[1].tok.Literal, 64)
			if err != nil {
				yylex.(*lexer).Error(err.Error())
			}

			yyVAL.expr = &ast.FloatLiteral{
				Token: yyDollar[1].tok,
				Value: f,
			}
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:328
		{
			if n := len(yyDollar[1].tok.Literal); n >= 2 {
				c, _, _, err := strconv.UnquoteChar(yyDollar[1].tok.Literal[1:n-1], '\'')
				if err != nil {
					yylex.(*lexer).Error(err.Error())
				} else {
					yyVAL.expr = &ast.CharLiteral{
						Token: yyDollar[1].tok,
						Value: c,
					}
				}
			}
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:342
		{
			yyVAL.expr = &ast.StringLiteral{
				Token: yyDollar[1].tok,
			}
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
//line rose.y:348
		{
			yyVAL.expr = &ast.RawStringLiteral{
				Token: yyDollar[1].tok,
			}
		}
	}
	goto yystack /* stack new state and value */
}
