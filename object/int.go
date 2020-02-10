package object

import (
	"strconv"
)

func (i Int) Type() ObjectType { return INTEGER_OBJ }
func (i Int) Truthy() bool     { return int64(i) != 0 }
func (i Int) Repr() string     { return strconv.FormatInt(int64(i), 10) }
