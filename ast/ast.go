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

type TypeName struct {
	Token token.Token
}

func (tn *TypeName) TokenLiteral() string { return tn.Token.Literal }
func (tn *TypeName) String() string       { return tn.Token.Literal }

//
// Statements
//

type VarDeclStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarDeclStatement) statementNode()       {}
func (vs *VarDeclStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarDeclStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	return out.String()
}

//
// Expressions
//

type Identifier struct {
	Token token.Token // the token.IDENT token
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Token.Literal }

type Nil struct {
	Token token.Token
}

func (n *Nil) expressionNode()      {}
func (n *Nil) TokenLiteral() string { return n.Token.Literal }
func (n *Nil) String() string       { return n.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

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
