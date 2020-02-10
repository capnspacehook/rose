package object

import (
	"strconv"
)

func (f Float) Type() ObjectType { return FLOAT_OBJ }
func (f Float) Truthy() bool     { return float64(f) > 0 }
func (f Float) Repr() string     { return strconv.FormatFloat(float64(f), 'g', -1, 64) }
