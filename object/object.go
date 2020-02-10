package object

import "github.com/capnspacehook/rose/token"

type ObjectType string

const (
	NIL_OBJ     ObjectType = "nil"
	BOOL_OBJ               = "bool"
	INTEGER_OBJ            = "int"
	FLOAT_OBJ              = "float"
	CHAR_OBJ               = "char"
	STRING_OBJ             = "string"
)

type Object interface {
	Type() ObjectType
	Truthy() bool
	Equals(rhs Object) bool
	String() string
}

type Nil struct{}

type Bool bool

type Int int64

type Float float64

type Char rune

type String string

type Nilable interface {
	IsNil() bool
}

type BinaryOperable interface {
	BinaryOp(op token.TokenType, rhs Object) (Object, error)
}

type Orderable interface {
	LessThan(rhs Object) bool
}
