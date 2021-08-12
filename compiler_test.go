package main

import (
	"strings"
	"testing"
)

func compareOps(t *testing.T, code string, expected ...Op) {
	parser := NewParser(strings.NewReader(code))
	compiler := NewCompiler(parser)

	actual := compiler.ReadExpressions()
	if len(actual) != len(expected) {
		t.Errorf("Expected %d ops, but got %d instead", len(expected), len(actual))
	}

	for i, op := range actual {
		if op != expected[i] {
			t.Errorf("Op %d differs: should be %v, but got %v instead.", i, expected[i], op)
		}
	}
}

func TestWordCompile(t *testing.T) {
	compareOps(t, ": foo 1 . ; foo",
			Op{OP_CALL, 0, Token{FUNCALL_TOKEN, 0, "foo"}},
	)
}

func assertPanic(t *testing.T, code string) {
	defer func() {
		if r := recover(); r == nil {
				t.Errorf("The code did not panic")
		}
	}()

	parser := NewParser(strings.NewReader(code))
	compiler := NewCompiler(parser)
	compiler.ReadExpressions()
}

func TestSpuriousSemicolon(t *testing.T) {
	assertPanic(t, "; foo 1 . ;")
}

func TestMissingSemicolon(t *testing.T) {
	assertPanic(t, ": foo 1 .")
}
