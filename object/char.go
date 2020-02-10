package object

import (
	"strconv"
)

func (c Char) Type() ObjectType { return CHAR_OBJ }
func (c Char) Truthy() bool     { return rune(c) != 0 }
func (c Char) Repr() string     { return strconv.QuoteRuneToGraphic(rune(c)) }
