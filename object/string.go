package object

func (s String) Type() ObjectType       { return STRING_OBJ }
func (s String) Truthy() bool           { return len(string(s)) > 0 }
func (s String) Equals(rhs Object) bool { return string(s) == string(rhs.(String)) }
func (s String) String() string         { return string(s) }
