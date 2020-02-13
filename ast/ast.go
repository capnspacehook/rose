// Copyright 2009 The Go Authors. All rights reserved.

// Package ast declares the types used to represent syntax trees for Rose
// packages.

package ast

import (
	"github.com/capnspacehook/rose/token"
)

// -----------------------------------------------------------------------------
// Interfaces
//
// There are 3 main classes of nodes: Expressions and type nodes,
// statement nodes, and declaration nodes. The node names usually
// match the corresponding Go spec production names to which they
// correspond. The node fields correspond to the individual parts
// of the respective productions.
//
// All nodes contain position information marking the beginning of
// the corresponding source text segment; it is accessible via the
// Pos accessor method. Nodes may contain additional position info
// for language constructs where comments may be found between parts
// of the construct (typically any larger, parenthesized subpart).
// That position information is needed to properly position comments
// when printing the construct.

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

// -----------------------------------------------------------------------------
// Expressions and types

// A BadExpr node is a placeholder for expressions containing
// syntax errors for which no correct expression nodes can be
// created.
type BadExpr struct {
	From, To token.Pos // position range of bad expression
}

// An Ident node represents an identifier.
type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name
	Obj     *Object   // denoted object; or nil
}

// A BasicLit node represents a literal of basic type.
type BasicLit struct {
	ValuePos token.Pos   // literal position
	Kind     token.Token // token.INT, token.FLOAT, token.CHAR, token.STRING, or token.RAW_STRING
	Value    string      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
}

// A ParenExpr node represents a parenthesized expression.
type ParenExpr struct {
	Lparen token.Pos // position of "("
	Expr   Expr      // parenthesized expression
	Rparen token.Pos // position of ")"
}

// A UnaryExpr node represents a unary expression.
type UnaryExpr struct {
	OpPos token.Pos   // position of Op
	Op    token.Token // operator
	Expr  Expr        // operand
}

// A BinaryExpr node represents a binary expression.
type BinaryExpr struct {
	Lhs   Expr        // left operand
	OpPos token.Pos   // position of Op
	Op    token.Token // operator
	Rhs   Expr        // right operand
}

// Pos and End implementations for expression/type nodes.
func (x *BadExpr) Pos() token.Pos    { return x.From }
func (x *Ident) Pos() token.Pos      { return x.NamePos }
func (x *BasicLit) Pos() token.Pos   { return x.ValuePos }
func (x *ParenExpr) Pos() token.Pos  { return x.Lparen }
func (x *UnaryExpr) Pos() token.Pos  { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos { return x.Lhs.Pos() }

func (x *BadExpr) End() token.Pos    { return x.To }
func (x *Ident) End() token.Pos      { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos   { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *ParenExpr) End() token.Pos  { return x.Rparen + 1 }
func (x *UnaryExpr) End() token.Pos  { return x.Expr.End() }
func (x *BinaryExpr) End() token.Pos { return x.Rhs.End() }

// exprNode() ensures that only expression/type nodes can be
// assigned to an Expr.
func (*BadExpr) exprNode()    {}
func (*Ident) exprNode()      {}
func (*BasicLit) exprNode()   {}
func (*ParenExpr) exprNode()  {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}

// -----------------------------------------------------------------------------
// Convenience functions for Idents

// IsExported reports whether name starts with an upper-case letter.
func IsExported(name string) bool { return token.IsExported(name) }

// IsExported reports whether id starts with an upper-case letter.
func (id *Ident) IsExported() bool { return token.IsExported(id.Name) }

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

// -----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.

// A BadStmt node is a placeholder for statements containing
// syntax errors for which no correct statement nodes can be
// created.
type BadStmt struct {
	From, To token.Pos // position range of bad statement
}

// A DeclStmt node represents a declaration in a statement list.
type DeclStmt struct {
	Decl Decl // *GenDecl with CONST, TYPE, or VAR token
}

// An EmptyStmt node represents an empty statement.
// The "position" of the empty statement is the position
// of the immediately following (explicit or implicit) semicolon.
type EmptyStmt struct {
	Semicolon token.Pos // position of following ";"
	Implicit  bool      // if set, ";" was omitted in the source
}

// An ExprStmt node represents a (stand-alone) expression
// in a statement list.
type ExprStmt struct {
	Expr Expr // expression
}

// An IncDecStmt node represents an increment or decrement statement.
type IncDecStmt struct {
	Expr   Expr
	TokPos token.Pos   // position of Tok
	Tok    token.Token // INC or DEC
}

// An AssignStmt node represents an assignment or
// a short variable declaration.
type AssignStmt struct {
	Lhs    []Expr
	TokPos token.Pos   // position of Tok
	Tok    token.Token // assignment token, DEFINE
	Rhs    []Expr
}

// Pos and End implementations for statement nodes.

func (s *BadStmt) Pos() token.Pos    { return s.From }
func (s *DeclStmt) Pos() token.Pos   { return s.Decl.Pos() }
func (s *EmptyStmt) Pos() token.Pos  { return s.Semicolon }
func (s *ExprStmt) Pos() token.Pos   { return s.Expr.Pos() }
func (s *IncDecStmt) Pos() token.Pos { return s.Expr.Pos() }
func (s *AssignStmt) Pos() token.Pos { return s.Lhs[0].Pos() }

func (s *BadStmt) End() token.Pos  { return s.To }
func (s *DeclStmt) End() token.Pos { return s.Decl.End() }
func (s *EmptyStmt) End() token.Pos {
	if s.Implicit {
		return s.Semicolon
	}
	return s.Semicolon + 1 /* len(";") */
}
func (s *ExprStmt) End() token.Pos { return s.Expr.End() }
func (s *IncDecStmt) End() token.Pos {
	return s.TokPos + 2 /* len("++") */
}
func (s *AssignStmt) End() token.Pos { return s.Rhs[len(s.Rhs)-1].End() }

// stmtNode() ensures that only statement nodes can be
// assigned to a Stmt.
func (*BadStmt) stmtNode()    {}
func (*DeclStmt) stmtNode()   {}
func (*EmptyStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()   {}
func (*IncDecStmt) stmtNode() {}
func (*AssignStmt) stmtNode() {}

// -----------------------------------------------------------------------------
// Declarations

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.

// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
type Spec interface {
	Node
	specNode()
}

// A ValueSpec node represents a constant or variable declaration
// (ConstSpec or VarSpec production).
type ValueSpec struct {
	//Doc     *CommentGroup // associated documentation; or nil
	Names  []*Ident // value names (len(Names) > 0)
	Type   Expr     // value type; or nil
	Values []Expr   // initial values; or nil
	//Comment *CommentGroup // line comments; or nil
}

// Pos and End implementations for spec nodes.

func (s *ValueSpec) Pos() token.Pos { return s.Names[0].Pos() }
func (s *ValueSpec) End() token.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}

// specNode() ensures that only spec nodes can be
// assigned to a Spec.
func (*ValueSpec) specNode() {}

// A declaration is represented by one of the following declaration nodes.

// A BadDecl node is a placeholder for declarations containing
// syntax errors for which no correct declaration nodes can be
// created.
type BadDecl struct {
	From, To token.Pos // position range of bad declaration
}

// A GenDecl node (generic declaration node) represents an import,
// constant, type or variable declaration. A valid Lparen position
// (Lparen.IsValid()) indicates a parenthesized declaration.
//
// Relationship between Tok value and Specs element type:
//
//	token.IMPORT  *ImportSpec
//	token.CONST   *ValueSpec
//	token.TYPE    *TypeSpec
//	token.VAR     *ValueSpec
type GenDecl struct {
	//Doc    *CommentGroup // associated documentation; or nil
	TokPos token.Pos   // position of Tok
	Tok    token.Token // IMPORT, CONST, TYPE, VAR
	Lparen token.Pos   // position of '(', if any
	Specs  []Spec
	Rparen token.Pos // position of ')', if any
}

// Pos and End implementations for declaration nodes.

func (d *BadDecl) Pos() token.Pos { return d.From }
func (d *GenDecl) Pos() token.Pos { return d.TokPos }

func (d *BadDecl) End() token.Pos { return d.To }
func (d *GenDecl) End() token.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Specs[0].End()
}

// declNode() ensures that only declaration nodes can be
// assigned to a Decl.
func (*BadDecl) declNode() {}
func (*GenDecl) declNode() {}

// ----------------------------------------------------------------------------
// Files and packages

// A File node represents a Go source file.
type File struct {
	Package token.Pos // position of "package" keyword
	Name    *Ident    // package name
	Stmts   []Stmt    // top-level statements; or nil
	Scope   *Scope    // package scope (this file only)
	//Imports    []*ImportSpec   // imports in this file
	Unresolved []*Ident // unresolved identifiers in this file
}

func (f *File) Pos() token.Pos { return f.Package }
func (f *File) End() token.Pos {
	if n := len(f.Stmts); n > 0 {
		return f.Stmts[n-1].End()
	}
	return f.Name.End()
}
