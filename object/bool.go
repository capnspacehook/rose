package object

func (b Bool) Type() ObjectType { return BOOL_OBJ }
func (b Bool) Truthy() bool     { return bool(b) }
func (b Bool) Repr() string {
	if bool(b) {
		return "true"
	}

	return "false"
}
