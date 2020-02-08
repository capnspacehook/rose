package parse

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	input := `var foo = 0x90
var bar = 56.0e1
var baz = 114_223_117
var assigned = foo
`

	statements, err := Parse(strings.NewReader(input))
	fmt.Println(statements, err)
}
