package object

func (n Nil) Type() ObjectType { return NIL_OBJ }
func (n Nil) Truthy() bool     { return false }
func (n Nil) Repr() string     { return "<nil>" }
