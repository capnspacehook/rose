package ast

import (
	"bytes"

	"github.com/capnspacehook/rose/token"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

//
// Abstract Nodes
//

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string { return "" }
func (p *Program) String() string {
	var buf bytes.Buffer

	for _, statement := range p.Statements {
		buf.WriteString(statement.String())
	}

	return buf.String()
}

type TypeName struct {
	Token token.Token
}

func (tn *TypeName) TokenLiteral() string { return tn.Token.Literal }
func (tn *TypeName) String() string       { return tn.Token.Literal }

//
// Statements
//

type ExprStatement struct {
	Expr Expression
}

func (es *ExprStatement) statementNode()       {}
func (es *ExprStatement) TokenLiteral() string { return es.Expr.TokenLiteral() }
func (es *ExprStatement) String() string       { return es.Expr.String() }

type VarDeclStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *TypeName
	Value Expression
}

func (vs *VarDeclStatement) statementNode()       {}
func (vs *VarDeclStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarDeclStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	if vs.Type != nil {
		out.WriteString(" " + vs.Type.String())
	}
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	return out.String()
}

type ConstDeclStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *TypeName
	Value Expression
}

func (cs *ConstDeclStatement) statementNode()       {}
func (cs *ConstDeclStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstDeclStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	if cs.Type != nil {
		out.WriteString(" " + cs.Type.String())
	}
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}

	return out.String()
}

type AssignmentStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" " + as.TokenLiteral() + " ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	return out.String()
}

//
// Expressions
//

type NilLiteral struct {
	Token token.Token
}

func (nl *NilLiteral) expressionNode()      {}
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NilLiteral) String() string       { return nl.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type CharLiteral struct {
	Token token.Token
	Value rune
}

func (cl *CharLiteral) expressionNode()      {}
func (cl *CharLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *CharLiteral) String() string       { return cl.Token.Literal }

type StringLiteral struct {
	Token token.Token
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type RawStringLiteral struct {
	Token token.Token
}

func (rl *RawStringLiteral) expressionNode()      {}
func (rl *RawStringLiteral) TokenLiteral() string { return rl.Token.Literal }
func (rl *RawStringLiteral) String() string       { return rl.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Token.Literal }

type ParenExpression struct {
	Lparen token.Token
	Expr   Expression
	Rparen token.Token
}

func (pe *ParenExpression) expressionNode()      {}
func (pe *ParenExpression) TokenLiteral() string { return pe.Lparen.Literal }
func (pe *ParenExpression) String() string {
	return "(" + pe.Expr.String() + ")"
}

type UnaryExpression struct {
	Token token.Token
	Value Expression
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) TokenLiteral() string { return ue.Token.Literal }
func (ue *UnaryExpression) String() string {
	return ue.Token.Literal + ue.Value.String()
}

type BinaryExpression struct {
	Lhs   Expression
	Token token.Token
	Rhs   Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BinaryExpression) String() string {
	return be.Lhs.String() + " " + be.Token.Literal + " " + be.Rhs.String()
}

type Conversion struct {
	Type  *TypeName
	Value Expression
}

func (c *Conversion) expressionNode()      {}
func (c *Conversion) TokenLiteral() string { return c.Type.TokenLiteral() }
func (c *Conversion) String() string {
	return c.Type.String() + "(" + c.Value.String() + ")"
}
