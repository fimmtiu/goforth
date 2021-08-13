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

func runCodeWithBuiltins(code string) {
	vm := NewVirtualMachine()
	compiler := NewCompiler(vm)
	compiler.LoadBuiltins()
	compiler.LoadCode(strings.NewReader(code))
	vm.Run()
}

func ExampleVirtualMachine_addition_and_printing() {
	runCode(": foo ( -- n ) 1 2 + ; foo .")
	// Output: 3
}


func ExampleVirtualMachine_modulus1() {
	runCode("3 2 mod .")
	// Output: 1
}

func ExampleVirtualMachine_modulus2() {
	runCode("6 2 mod .")
	// Output: 0
}

func ExampleVirtualMachine_dup() {
	runCode("13 dup . .")
	// Output: 1313
}

func ExampleVirtualMachine_if_then_true() {
	runCode("31337 1 if . then")
	// Output: 31337
}

func ExampleVirtualMachine_cr() {
	runCodeWithBuiltins("1 . cr 2 .")
	// Output:
	// 1
	// 2
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
