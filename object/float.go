package object

import (
	"strconv"
)

func (f Float) Type() ObjectType       { return FLOAT_OBJ }
func (f Float) Truthy() bool           { return float64(f) != 0 }
func (f Float) Equals(rhs Object) bool { return float64(f) == float64(rhs.(Float)) }
func (f Float) String() string         { return strconv.FormatFloat(float64(f), 'g', -1, 64) }
