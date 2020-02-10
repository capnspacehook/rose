package eval

import (
	"github.com/capnspacehook/rose/ast"
	"github.com/capnspacehook/rose/object"
)

var (
	NIL   *object.Nil
	TRUE  *object.Bool
	FALSE *object.Bool
)

func init() {
	NIL = &object.Nil{}
	t := object.Bool(true)
	f := object.Bool(false)

	TRUE = &t
	FALSE = &f
}

func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
		// Statements
	case *ast.ExprStatement:
		return Eval(node.Expr)

		// Expressions
	case *ast.NilLiteral:
		return NIL, nil
	case *ast.BooleanLiteral:
		return nativeBoolToBoolObj(node.Value), nil
	case *ast.IntegerLiteral:
		return object.Int(node.Value), nil
	case *ast.FloatLiteral:
		return object.Float(node.Value), nil
	case *ast.CharLiteral:
		return object.Char(node.Value), nil
	case *ast.StringLiteral:
		return object.String(node.Token.Literal), nil
	}

	return nil, nil
}

func evalStatements(stmts []ast.Statement) (obj object.Object, err error) {
	for _, statement := range stmts {
		obj, err = Eval(statement)
		if err != nil {
			return nil, err
		}
	}

	return
}

func nativeBoolToBoolObj(input bool) *object.Bool {
	if input {
		return TRUE
	}
	return FALSE
}
