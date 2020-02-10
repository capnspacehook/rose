package object

func (b Bool) Type() ObjectType       { return BOOL_OBJ }
func (b Bool) Truthy() bool           { return bool(b) }
func (b Bool) Equals(rhs Object) bool { return bool(b) == bool(rhs.(Bool)) }
func (b Bool) String() string {
	if bool(b) {
		return "true"
	}

	return "false"
}
