package object

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
	Repr() string
}

type Adder interface {
	Add(rhs Object) Object
}

type Nil struct{}

type Bool bool

type Int int64

type Float float64

type Char rune

type String string
