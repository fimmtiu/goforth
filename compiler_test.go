package main

import (
	"strings"
	"testing"
)

func compareOps(t *testing.T, code string, expected ...AbstractOp) *Compiler {
	c := NewCompiler(NewVirtualMachine())
	c.parser = NewParser(strings.NewReader(code))

	actual := c.Compile()
	if len(actual) != len(expected) {
		t.Errorf("Expected %d ops, but got %d instead", len(expected), len(actual))
	}

	for i, op := range actual {
		if op != expected[i] {
			t.Errorf("AbstractOp %d differs: should be %v, but got %v instead.", i, expected[i], op)
		}
	}

	return c
}

func wordsEqual(one Word, two Word) bool {
	if one.Name != two.Name || len(one.Ops) != len(two.Ops) {
		return false
	}

	for i, op := range one.Ops {
		if op != two.Ops[i] {
			return false
		}
	}

	return true
}

func assertPanic(t *testing.T, code string) {
	defer func() {
		if r := recover(); r == nil {
				t.Errorf("Expected a panic, but didn't get one: %s", code)
		}
	}()

	c := NewCompiler(NewVirtualMachine())
	c.LoadCode(strings.NewReader(code))
}

func assertPackedOpsEqual(t *testing.T, actual []PackedOp, expected []PackedOp) {
	if len(actual) != len(expected) {
		t.Errorf("Size mismatch for packed ops: expected %d but got %d", len(expected), len(actual))
	}

	for i, op := range actual {
		if op != expected[i] {
			t.Errorf("Op %d was expected to be %08x, but was %08x", i, expected[i], op)
		}
	}
}

func TestWordCompile(t *testing.T) {
	c := compareOps(t, ": foo 1 . ; foo",
			AbstractOp{OP_CALL, 0, StringDatum{"foo"}},
	)

	if len(c.words) != 1 {
		t.Errorf("Expected 1 word, but got %d", len(c.words))
	}

	foo := Word{"foo", []AbstractOp{{OP_PUSH, 0, IntegerDatum{1}}, {OP_PRINT, 0, VoidDatum{}}, {OP_RETURN, 0, VoidDatum{}}}}

	if !wordsEqual(c.words[0], foo) {
		t.Errorf("Expected newly defined word to be %v, but got %v", foo, c.words[0])
	}
}

func TestWordOpPacking(t *testing.T) {
	c := NewCompiler(NewVirtualMachine())
	c.LoadCode(strings.NewReader(": foo 1 2 + ; foo ."))

	if len(c.vm.Dict) != 2 {
		t.Errorf("Expected 2 entries in the dictionary, but got %d.", len(c.vm.Dict))
	}
	if c.vm.Dict["foo"] != 0 {
		t.Errorf("Expected foo to start at offset 0, but it's at %d.", c.vm.Dict["foo"])
	}
	if c.vm.Dict["top-level code"] != 4 {
		t.Errorf("Expected top-level code to start at offset 4, but it's at %d.", c.vm.Dict["top-level code"])
	}

	assertPackedOpsEqual(t, c.vm.Code, []PackedOp{
		0x00000002, // OP_PUSH 1  [start of foo]
		0x00000102, // OP_PUSH 2
		0x00000007, // OP_ADD
		0x00000001, // OP_RETURN
		0x00000003, // OP_CALL 0  [start of top-level code]
		0x00000006, // OP_PRINT
		0x00000001, // OP_RETURN
	})

	if c.vm.Ip != 4 {
		t.Errorf("Expected instruction pointer to be 4, but got %d.", len(c.vm.Dict))
	}
}

func TestIfCompile(t *testing.T) {
	compareOps(t, "1 if 2 then",
			AbstractOp{OP_PUSH, 0, IntegerDatum{1}},
			AbstractOp{OP_JUMP_IF_NOT, 2, VoidDatum{}},
			AbstractOp{OP_PUSH, 0, IntegerDatum{2}},
	)
}

func TestIfElseCompile(t *testing.T) {
	compareOps(t, "1 if 2 else 3 then",
			AbstractOp{OP_PUSH, 0, IntegerDatum{1}},
			AbstractOp{OP_JUMP_IF_NOT, 3, VoidDatum{}},
			AbstractOp{OP_PUSH, 0, IntegerDatum{2}},
			AbstractOp{OP_JUMP, 2, VoidDatum{}},
			AbstractOp{OP_PUSH, 0, IntegerDatum{3}},
	)
}

func TestIfOpPacking(t *testing.T) {
	c := NewCompiler(NewVirtualMachine())
	c.LoadCode(strings.NewReader("1 if 2 then"))

	assertPackedOpsEqual(t, c.vm.Code, []PackedOp{
		0x00000002, // OP_PUSH 1
		0x00000305, // OP_JUMP_IF_NOT 3
		0x00000102, // OP_PUSH 2
		0x00000001, // OP_RETURN
	})

	if c.vm.Ip != 0 {
		t.Errorf("Expected instruction pointer to be 0, but got %d.", len(c.vm.Dict))
	}
}

func TestIfElseOpPacking(t *testing.T) {
	c := NewCompiler(NewVirtualMachine())
	c.LoadCode(strings.NewReader("1 if 2 else 3 then"))

	assertPackedOpsEqual(t, c.vm.Code, []PackedOp{
		0x00000002, // OP_PUSH 1
		0x00000405, // OP_JUMP_IF_NOT 4
		0x00000102, // OP_PUSH 2
		0x00000504, // OP_JUMP 5
		0x00000202, // OP_PUSH 3
		0x00000001, // OP_RETURN
	})
}

func TestStoreFetchOpPacking(t *testing.T) {
	c := NewCompiler(NewVirtualMachine())
	c.LoadCode(strings.NewReader("1 foo ! foo @"))
	c.vm.printDisassembly()

	assertPackedOpsEqual(t, c.vm.Code, []PackedOp{
		0x00000002, // OP_PUSH 1
		0x0000010c, // OP_STORE foo
		0x0000020d, // OP_FETCH foo
		0x00000001, // OP_RETURN
	})
}

func TestSpuriousSemicolon(t *testing.T) {
	assertPanic(t, "; foo 1 . ;")
}

func TestMissingSemicolon(t *testing.T) {
	assertPanic(t, ": foo 1 .")
}

func TestSpuriousElse(t *testing.T) {
	assertPanic(t, "1 else .")
}

func TestSpuriousThen(t *testing.T) {
	assertPanic(t, "1 then .")
}

func TestUnterminatedIf(t *testing.T) {
	assertPanic(t, "1 if foo")
	assertPanic(t, "1 if foo else bar")
}

func TestCompileAddition(t *testing.T) {
	compareOps(t, "foo ( n1 n2 -- n' ) 1 2 + .",
			AbstractOp{OP_CALL, 0, StringDatum{"foo"}},
			AbstractOp{OP_PUSH, 0, IntegerDatum{1}},
			AbstractOp{OP_PUSH, 0, IntegerDatum{2}},
			AbstractOp{OP_ADD, 0, VoidDatum{}},
			AbstractOp{OP_PRINT, 0, VoidDatum{}},
	)
}

func TestCompilePrintString(t *testing.T) {
	compareOps(t, `"foo" .`,
			AbstractOp{OP_PUSH, 0, StringDatum{"foo"}},
			AbstractOp{OP_PRINT, 0, VoidDatum{}},
	)
}
