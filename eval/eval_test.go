package eval

import (
	"fmt"
	"strings"
	"testing"

	"github.com/capnspacehook/rose/parse"
)

func TestEval(t *testing.T) {
	input := `'☺'`

	program, err := parse.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	obj, err := Eval(program)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(obj)
}
