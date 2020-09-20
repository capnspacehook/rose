package parser

import (
	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/token"
)

func (p *Parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case token.CONST, token.LET, token.VAR:
		s = &ast.DeclStmt{Decl: p.parseDecl()}
	default:
		// no statement found
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.advance(stmtStart)
		s = &ast.BadStmt{From: pos, To: p.pos}
	}

	return
}

type parseSpecFunc func(keyword token.Token, i int) ast.Spec

func (p *Parser) parseDecl() ast.Decl {
	if p.trace {
		defer un(trace(p, "Declaration"))
	}

	var f parseSpecFunc
	switch p.tok {
	case token.CONST, token.LET, token.VAR:
		f = p.parseValueSpec
	default:
		pos := p.pos
		p.errorExpected(pos, "declaration")
		p.advance(stmtStart)
		return &ast.BadDecl{From: pos, To: p.pos}
	}

	return p.parseGenDecl(p.tok, f)

}

func (p *Parser) parseGenDecl(keyword token.Token, f parseSpecFunc) *ast.GenDecl {
	if p.trace {
		defer un(trace(p, "GenDecl("+keyword.String()+")"))
	}

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
	if p.trace {
		defer un(trace(p, keyword.String()+"Spec"))
	}

	pos := p.pos
	idents := p.parseIdentList()
	typ := p.tryType()
	var values []ast.Expr
	if p.tok == token.ASSIGN {
		if keyword == token.VAR {
			p.error(pos, "initialization is not allowed in a 'var' statement")
		}
		p.next()
		values = p.parseRhsList()
	}
	p.expectSemi()

	switch keyword {
	case token.LET:
		if values == nil {
			p.error(pos, "missing initialization")
		}
	case token.VAR:
		if typ == nil {
			p.error(pos, "missing variable type")
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
