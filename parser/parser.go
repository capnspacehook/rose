package parser

import (
	"fmt"
	"io"

	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/lexer"
	"github.com/capnspacehook/rose/token"
)

var typeNames = map[string]bool{
	"any":    true,
	"bool":   true,
	"int":    true,
	"float":  true,
	"string": true,
}

var boolConsts = map[string]bool{
	"false": false,
	"true":  true,
}

type Parser struct {
	fset   *token.FileSet
	errors lexer.ErrorList
	lexer  lexer.Lexer

	file *token.File

	// Next token
	pos token.Pos   // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

	// Error recovery
	// (used to limit the number of calls to parser.advance
	// w/o making scanning progress - avoids potential endless
	// loops across multiple parser functions during error recovery)
	syncPos token.Pos // last synchronization position
	syncCnt int       // number of parser.advance calls without progress

	// Non-syntactic parser control
	exprLev int  // < 0: in control clause, >= 0: in expression
	inRhs   bool // if set, the parser is parsing a rhs expression

	// Ordinary identifier scopes
	pkgScope   *ast.Scope   // pkgScope.Outer == nil
	topScope   *ast.Scope   // top-most scope; may be pkgScope
	unresolved []*ast.Ident // unresolved identifiers
}

func NewParser(fset *token.FileSet) (p Parser) {
	p.fset = fset
	return
}

// ----------------------------------------------------------------------------
// Scoping support

func (p *Parser) openScope() {
	p.topScope = ast.NewScope(p.topScope)
}

func (p *Parser) closeScope() {
	p.topScope = p.topScope.Outer
}

// TODO: (capnspacehook) make decl's type ast.Node?
func (p *Parser) declare(decl, data interface{}, scope *ast.Scope, kind ast.ObjKind, idents ...*ast.Ident) {
	for _, ident := range idents {
		assert(ident.Obj == nil, "identifier already declared or resolved")
		obj := ast.NewObj(kind, ident.Name)
		// remember the corresponding declaration for redeclaration
		// errors and global variable resolution/typechecking phase
		obj.Decl = decl
		obj.Data = data
		ident.Obj = obj
		if ident.Name != "_" {
			if alt := scope.Insert(obj); alt != nil /*&& p.mode&DeclarationErrors != 0*/ {
				prevDecl := ""
				if pos := alt.Pos(); pos.IsValid() {
					prevDecl = fmt.Sprintf("\n\tprevious declaration at %s", p.file.Position(pos))
				}
				p.error(ident.Pos(), fmt.Sprintf("%s redeclared in this block%s", ident.Name, prevDecl))
			}
		}
	}
}

// The unresolved object is a sentinel to mark identifiers that have been added
// to the list of unresolved identifiers. The sentinel is only used for verifying
// internal consistency.
var unresolved = new(ast.Object)

// If x is an identifier, tryResolve attempts to resolve x by looking up
// the object it denotes. If no object is found and collectUnresolved is
// set, x is marked as unresolved and collected in the list of unresolved
// identifiers.
func (p *Parser) tryResolve(x ast.Expr, collectUnresolved bool) {
	// nothing to do if x is not an identifier or the blank identifier
	ident, _ := x.(*ast.Ident)
	if ident == nil {
		return
	}
	//assert(ident.Obj == nil, "identifier already declared or resolved")
	if ident.Name == "_" {
		return
	}

	// try to resolve the identifier
	for s := p.topScope; s != nil; s = s.Outer {
		if obj := s.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return
		}
	}
	// all local scopes are known, so any unresolved identifier
	// must be found either in the file scope, package scope
	// (perhaps in another file), or universe scope --- collect
	// them so that they can be resolved later
	if collectUnresolved {
		ident.Obj = unresolved
		p.unresolved = append(p.unresolved, ident)
	}
}

func (p *Parser) resolve(x ast.Expr) {
	p.tryResolve(x, true)
}

// ----------------------------------------------------------------------------
// Error handling

// A bailout panic is raised to indicate early termination.
type bailout struct{}

//TODO: (capnspacehook) Fix token.Position/scanner.Position conflict
func (p *Parser) error(pos token.Pos, msg string) {
	//epos := p.file.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	/*if p.mode&AllErrors == 0 {
		n := len(p.errors)
		if n > 0 && p.errors[n-1].Pos.Line == epos.Line {
			return // discard - likely a spurious error
		}
		if n > 10 {
			panic(bailout{})
		}
	}*/

	//p.errors.Add(epos, msg)
}

func (p *Parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			msg += ", found newline"
		case p.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}

// ----------------------------------------------------------------------------
// Parsing helpers

func (p *Parser) next() {
	p.pos, p.tok, p.lit = p.lexer.Lex()
}

func (p *Parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress

	return pos
}

func (p *Parser) expectSemi() {
	// semicolon is optional before a closing ')' or '}'
	if p.tok != token.RPAREN && p.tok != token.RBRACE {
		if p.tok == token.SEMICOLON {
			p.next()
		} else {
			p.errorExpected(p.pos, "';'")
			//p.advance(stmtStart)
		}
	}
}

func assert(cond bool, msg string) {
	if !cond {
		panic("rose/parser internal error: " + msg)
	}
}

// advance consumes tokens until the current token p.tok
// is in the 'to' set, or token.EOF. For error recovery.
func (p *Parser) advance(to map[token.Token]bool) {
	for ; p.tok != token.EOF; p.next() {
		if to[p.tok] {
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 advance calls without
			// progress. Otherwise consume at least one token to
			// avoid an endless parser loop (it is possible that
			// both parseOperand and parseStmt call advance and
			// correctly do not advance, thus the need for the
			// invocation limit p.syncCnt).
			if p.pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
				return
			}
			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCnt = 0
				return
			}
			// Reaching here indicates a parser bug, likely an
			// incorrect token list in this function, but it only
			// leads to skipping of possibly correct code if a
			// previous error is present, and thus is preferred
			// over a non-terminating parse.
		}
	}
}

var stmtStart = map[token.Token]bool{
	/*token.BREAK:       true,
	token.CONST:       true,
	token.CONTINUE:    true,
	token.DEFER:       true,
	token.FALLTHROUGH: true,
	token.FOR:         true,
	token.GO:          true,
	token.GOTO:        true,
	token.IF:          true,
	token.RETURN:      true,
	token.SELECT:      true,
	token.SWITCH:      true,
	token.TYPE:        true,*/
	token.VAR: true,
}

var declStart = map[token.Token]bool{
	token.CONST: true,
	//token.TYPE:  true,
	token.VAR: true,
}

var exprEnd = map[token.Token]bool{
	token.COMMA:     true,
	token.COLON:     true,
	token.SEMICOLON: true,
	token.RPAREN:    true,
	token.RBRACK:    true,
	token.RBRACE:    true,
}

// safePos returns a valid file position for a given position: If pos
// is valid to begin with, safePos returns pos. If pos is out-of-range,
// safePos returns the EOF position.
//
// This is hack to work around "artificial" end positions in the AST which
// are computed by adding 1 to (presumably valid) token positions. If the
// token positions are invalid due to parse errors, the resulting end position
// may be past the file's EOF position, which would lead to panics if used
// later on.
func (p *Parser) safePos(pos token.Pos) (res token.Pos) {
	defer func() {
		if recover() != nil {
			res = token.Pos(p.file.Base() + p.file.Size()) // EOF position
		}
	}()
	_ = p.file.Offset(pos) // trigger a panic if position is out-of-range
	return pos
}

// ----------------------------------------------------------------------------
// Identifiers

func (p *Parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(token.IDENT)
	}

	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *Parser) parseIdentList() (list []*ast.Ident) {
	list = append(list, p.parseIdent())
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.parseIdent())
	}

	return
}

// ----------------------------------------------------------------------------
// Common productions

// If lhs is set and the result is an identifier, it is not resolved.
// The result may be a type or even a raw type ([...]int). Callers must
// check the result (using checkExpr or checkExprOrType), depending on
// context.
func (p *Parser) parseExpr(lhs bool) ast.Expr {
	return p.parseBinaryExpr(lhs, token.LowestPrec+1)
}

// If lhs is set, result list elements which are identifiers are not resolved.
func (p *Parser) parseExprList(lhs bool) (list []ast.Expr) {
	list = append(list, p.checkExpr(p.parseExpr(lhs)))
	for p.tok == token.COMMA {
		p.next()
		list = append(list, p.checkExpr(p.parseExpr(lhs)))
	}

	return
}

func (p *Parser) parseLhsList() []ast.Expr {
	old := p.inRhs
	p.inRhs = false
	list := p.parseExprList(true)
	switch p.tok {
	case token.COLON:
		// lhs of a label declaration or a communication clause of a select
		// statement (parseLhsList is not called when parsing the case clause
		// of a switch statement):
		// - labels are declared by the caller of parseLhsList
		// - for communication clauses, if there is a stand-alone identifier
		//   followed by a colon, we have a syntax error; there is no need
		//   to resolve the identifier in that case
	default:
		// identifiers must be declared elsewhere
		for _, x := range list {
			p.resolve(x)
		}
	}
	p.inRhs = old

	return list
}

func (p *Parser) parseRhsList() []ast.Expr {
	old := p.inRhs
	p.inRhs = true
	list := p.parseExprList(false)
	p.inRhs = old
	return list
}

// ----------------------------------------------------------------------------
// Types

func (p *Parser) parseType() ast.Expr {
	typ := p.tryType()

	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.advance(exprEnd)
		return &ast.BadExpr{From: pos, To: p.pos}
	}

	return typ
}

// If the result is an identifier, it is not resolved.
func (p *Parser) parseTypeName() ast.Expr {
	ident := p.parseIdent()
	// don't resolve ident yet - it may be a parameter or field name

	/*if p.tok == token.PERIOD {
		// ident is a package name
		p.next()
		p.resolve(ident)
		sel := p.parseIdent()
		return &ast.SelectorExpr{X: ident, Sel: sel}
	}*/

	return ident
}

// If the result is an identifier, it is not resolved.
func (p *Parser) tryIdentOrType() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.parseTypeName()
	/*case token.LBRACK:
		return p.parseArrayType()
	case token.STRUCT:
		return p.parseStructType()
	case token.FUNC:
		typ, _ := p.parseFuncType()
		return typ
	case token.INTERFACE:
		return p.parseInterfaceType()
	case token.MAP:
		return p.parseMapType()
	case token.CHAN, token.ARROW:
		return p.parseChanType()*/
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, Expr: typ, Rparen: rparen}
	}

	// no type found
	return nil
}

func (p *Parser) tryType() ast.Expr {
	typ := p.tryIdentOrType()
	if typ != nil {
		p.resolve(typ)
	}
	return typ
}

// ----------------------------------------------------------------------------
// Expressions

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseUnaryExpr(lhs bool) ast.Expr {
	switch p.tok {
	case token.ADD, token.SUB, token.NOT, token.INVT, token.AND:
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr(false)
		return &ast.UnaryExpr{OpPos: pos, Op: op, Expr: p.checkExpr(x)}

		/*case token.ARROW:
		// channel type or receive expression
		arrow := p.pos
		p.next()

		// If the next token is token.CHAN we still don't know if it
		// is a channel type or a receive operation - we only know
		// once we have found the end of the unary expression. There
		// are two cases:
		//
		//   <- type  => (<-type) must be channel type
		//   <- expr  => <-(expr) is a receive from an expression
		//
		// In the first case, the arrow must be re-associated with
		// the channel type parsed already:
		//
		//   <- (chan type)    =>  (<-chan type)
		//   <- (chan<- type)  =>  (<-chan (<-type))

		x := p.parseUnaryExpr(false)

		// determine which case we have
		if typ, ok := x.(*ast.ChanType); ok {
			// (<-type)

			// re-associate position info and <-
			dir := ast.SEND
			for ok && dir == ast.SEND {
				if typ.Dir == ast.RECV {
					// error: (<-type) is (<-(<-chan T))
					p.errorExpected(typ.Arrow, "'chan'")
				}
				arrow, typ.Begin, typ.Arrow = typ.Arrow, arrow, arrow
				dir, typ.Dir = typ.Dir, ast.RECV
				typ, ok = typ.Value.(*ast.ChanType)
			}
			if dir == ast.SEND {
				p.errorExpected(arrow, "channel type")
			}

			return x
		}

		// <-(expr)
		return &ast.UnaryExpr{OpPos: arrow, Op: token.ARROW, Expr: p.checkExpr(x)}*/
	}

	return p.parsePrimaryExpr(lhs)
}

// parseOperand may return an expression or a raw type (incl. array
// types of the form [...]T. Callers must verify the result.
// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseOperand(lhs bool) ast.Expr {
	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if !lhs {
			p.resolve(x)
		}
		return x

	case token.INT, token.FLOAT, token.CHAR, token.STRING, token.RAW_STRING:
		x := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x

	case token.LPAREN:
		lparen := p.pos
		p.next()
		p.exprLev++
		x := p.parseRhsOrType() // types may be parenthesized: (some type)
		p.exprLev--
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, Expr: x, Rparen: rparen}

		/*case token.FUNC:
		return p.parseFuncTypeOrLit()*/
	}

	/*if typ := p.tryIdentOrType(); typ != nil {
		// could be type for composite literal or conversion
		_, isIdent := typ.(*ast.Ident)
		assert(!isIdent, "type cannot be identifier")
		return typ
	}*/

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	p.advance(stmtStart)
	return &ast.BadExpr{From: pos, To: p.pos}
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parsePrimaryExpr(lhs bool) ast.Expr {
	x := p.parseOperand(lhs)
	/*L:
	for {
		switch p.tok {
		case token.PERIOD:
		p.next()
		if lhs {
			p.resolve(x)
		}
		switch p.tok {
		case token.IDENT:
			x = p.parseSelector(p.checkExprOrType(x))
		case token.LPAREN:
			x = p.parseTypeAssertion(p.checkExpr(x))
		default:
			pos := p.pos
			p.errorExpected(pos, "selector or type assertion")
			p.next() // make progress
			sel := &ast.Ident{NamePos: pos, Name: "_"}
			x = &ast.SelectorExpr{X: x, Sel: sel}
		}
		case token.LBRACK:
			if lhs {
				p.resolve(x)
			}
			x = p.parseIndexOrSlice(p.checkExpr(x))
		case token.LPAREN:
			if lhs {
				p.resolve(x)
			}
			x = p.parseCallOrConversion(p.checkExprOrType(x))
		case token.LBRACE:
			if isLiteralType(x) && (p.exprLev >= 0 || !isTypeName(x)) {
				if lhs {
					p.resolve(x)
				}
				x = p.parseLiteralValue(x)
			} else {
				break L
			}
		default:
			break L
		}
		lhs = false // no need to try to resolve again
	}*/

	return x
}

func (p *Parser) tokPrec() (token.Token, int) {
	tok := p.tok
	if p.inRhs && tok == token.ASSIGN {
		tok = token.EQL
	}
	return tok, tok.Precedence()
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseBinaryExpr(lhs bool, prec1 int) ast.Expr {
	x := p.parseUnaryExpr(lhs)
	for {
		op, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		pos := p.expect(op)
		if lhs {
			p.resolve(x)
			lhs = false
		}
		y := p.parseBinaryExpr(false, oprec+1)
		x = &ast.BinaryExpr{Lhs: p.checkExpr(x), OpPos: pos, Op: op, Rhs: p.checkExpr(y)}
	}
}

func (p *Parser) parseRhsOrType() ast.Expr {
	old := p.inRhs
	p.inRhs = true
	x := p.checkExprOrType(p.parseExpr(false))
	p.inRhs = old
	return x
}

// If x is of the form (T), unparen returns unparen(T), otherwise it returns x.
func unparen(x ast.Expr) ast.Expr {
	if p, isParen := x.(*ast.ParenExpr); isParen {
		x = unparen(p.Expr)
	}
	return x
}

// checkExprOrType checks that x is an expression or a type
// (and not a raw type such as [...]T).
//
func (p *Parser) checkExprOrType(x ast.Expr) ast.Expr {
	switch unparen(x).(type) {
	case *ast.ParenExpr:
		panic("unreachable")
	case *ast.UnaryExpr:
		/*case *ast.ArrayType:
		if len, isEllipsis := t.Len.(*ast.Ellipsis); isEllipsis {
			p.error(len.Pos(), "expected array length, found '...'")
			x = &ast.BadExpr{From: x.Pos(), To: p.safePos(x.End())}
		}
		*/
	}

	// all other nodes are expressions or types
	return x
}

// checkExpr checks that x is an expression (and not a type).
func (p *Parser) checkExpr(x ast.Expr) ast.Expr {
	switch unparen(x).(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.BasicLit:
	//case *ast.FuncLit:
	//case *ast.CompositeLit:
	case *ast.ParenExpr:
		panic("unreachable")
	/*case *ast.SelectorExpr:
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
		// If t.Type == nil we have a type assertion of the form
		// y.(type), which is only allowed in type switch expressions.
		// It's hard to exclude those but for the case where we are in
		// a type switch. Instead be lenient and test this in the type
		// checker.
	case *ast.CallExpr:
	case *ast.StarExpr:*/
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
	default:
		// all other nodes are not proper expressions
		p.errorExpected(x.Pos(), "expression")
		x = &ast.BadExpr{From: x.Pos(), To: p.safePos(x.End())}
	}
	return x
}

// ----------------------------------------------------------------------------
// Parsing logic

func (p *Parser) ParseFile(filename string, src io.Reader, srcLen int) *ast.File {
	p.file = p.fset.AddFile(filename, -1, srcLen)
	p.lexer.Init(p.file, src, srcLen)

	p.next()

	p.openScope()
	var stmts []ast.Stmt
	for p.tok != token.EOF {
		stmts = append(stmts, p.parseStmt())
	}
	p.closeScope()

	return &ast.File{
		Stmts:      stmts,
		Unresolved: p.unresolved,
	}
}

func (p *Parser) parseStmt() (s ast.Stmt) {
	switch p.tok {
	case token.CONST, token.VAR:
		s = &ast.DeclStmt{Decl: p.parseDecl()}
	}

	return
}

type parseSpecFunc func(keyword token.Token, i int) ast.Spec

func (p *Parser) parseDecl() ast.Decl {
	var f parseSpecFunc
	switch p.tok {
	case token.CONST, token.VAR:
		f = p.parseValueSpec

	}

	return p.parseGenDecl(p.tok, f)

}

func (p *Parser) parseGenDecl(keyword token.Token, f parseSpecFunc) *ast.GenDecl {
	pos := p.expect(keyword)
	var lparen, rparen token.Pos
	var list []ast.Spec

	if p.tok == token.LPAREN {
		lparen = p.pos
		p.next()
		for i := 0; p.tok != token.RPAREN && p.tok != token.EOF; i++ {
			list = append(list, f(keyword, i))
		}

		rparen = p.expect(token.RPAREN)
		p.expectSemi()
	} else {
		list = append(list, f(keyword, 0))
	}

	return &ast.GenDecl{
		TokPos: pos,
		Tok:    keyword,
		Lparen: lparen,
		Specs:  list,
		Rparen: rparen,
	}
}

func (p *Parser) parseValueSpec(keyword token.Token, i int) ast.Spec {
	pos := p.pos
	idents := p.parseIdentList()
	typ := p.tryType()
	var values []ast.Expr
	// always permit optional initialization for more tolerant parsing
	if p.tok == token.ASSIGN {
		p.next()
		values = p.parseRhsList()
	}
	p.expectSemi()

	switch keyword {
	case token.VAR:
		if typ == nil && values == nil {
			p.error(pos, "missing variable type or initialization")
		}
	case token.CONST:
		if values == nil && (i == 0 || typ != nil) {
			p.error(pos, "missing constant value")
		}
	}

	spec := &ast.ValueSpec{
		Names:  idents,
		Type:   typ,
		Values: values,
	}
	kind := ast.Con
	if keyword == token.VAR {
		kind = ast.Var
	}
	p.declare(spec, i, p.topScope, kind, idents...)

	return spec
}
