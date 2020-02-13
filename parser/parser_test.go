package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/capnspacehook/rose/token"
)

var input = `var foo int = (5 ** ~5) & (3.4 ^ (78 % 5))`

func TestParser(t *testing.T) {
	p := NewParser(token.NewFileSet())
	astFile := p.ParseFile("test_file.rose", strings.NewReader(input), len(input))
	fmt.Println(astFile)
}
