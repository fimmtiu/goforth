package main

import (
	"strings"
)

func runCode(code string) {
	vm := NewVirtualMachine()
	compiler := NewCompiler(vm)
	compiler.LoadCode(strings.NewReader(code))
	vm.Run()
}

func ExampleVirtualMachine_addition_and_printing() {
	runCode(": foo ( -- n ) 1 2 + ; foo .")
	// Output: 3
}

func ExampleVirtualMachine_if_then_true() {
	runCode("31337 1 if . then")
	// Output: 31337
}

func ExampleVirtualMachine_if_then_false() {
	runCode("31337 0 if . then")
	// Output:
}

func ExampleVirtualMachine_if_else_then_true() {
	runCode("1 if 31337 else 69105 then .")
	// Output: 31337
}

func ExampleVirtualMachine_if_else_then_false() {
	runCode("0 if 31337 else 69105 then .")
	// Output: 69105
}
