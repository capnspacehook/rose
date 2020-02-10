package object

func (n Nil) Type() ObjectType       { return NIL_OBJ }
func (n Nil) Truthy() bool           { return false }
func (n Nil) Equals(rhs Object) bool { return rhs.(Nilable).IsNil() }
func (n Nil) String() string         { return "<nil>" }
