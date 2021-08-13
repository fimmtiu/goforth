package main

import (
	"bufio"
	"os"
)

func main() {
	vm := NewVirtualMachine()
	compiler := NewCompiler(vm)

	compiler.LoadBuiltins()
	compiler.LoadCode(bufio.NewReader(os.Stdin))
	vm.Run()
}
