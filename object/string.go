package object

func (s String) Type() ObjectType { return STRING_OBJ }
func (s String) Truthy() bool     { return len(string(s)) > 0 }
func (s String) Repr() string     { return string(s) }
