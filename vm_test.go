package main

import (
	"strings"
)

func ExampleVirtualMachine_basic_addition() {
	vm := NewVirtualMachine()
	compiler := NewCompiler(vm)
	compiler.LoadCode(strings.NewReader(": foo ( -- n ) 1 2 + ; foo ."))
	vm.Run()

	// Output: 3
}
