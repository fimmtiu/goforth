package main

import (
	"fmt"
)

type VirtualMachine struct {
	Heap []Datum
	Dict map[string]uint32
	Code []PackedOp
	Ip uint32

	dataStack []Datum
	callStack []uint32
	variables map[string]Datum
}

func NewVirtualMachine() *VirtualMachine {
	var vm VirtualMachine
	vm.Dict = make(map[string]uint32)
	vm.variables = make(map[string]Datum)
	return &vm
}

func (vm *VirtualMachine) Run() {
	// vm.printDisassembly()

	for {
		instruction := vm.Code[vm.Ip]
		opcode := uint8(instruction & 0xFF)
		arg := uint32(instruction >> 8)
		// fmt.Printf("Executing: [%d] opcode %d, arg %d\n", vm.Ip, opcode, arg)

		switch opcode {
		case OP_PRINT:
			printDatum(vm.popDataStack(), false)
		case OP_ADD:
			result := addNumbers(vm.popDataStack(), vm.popDataStack())
			vm.pushDataStack(result)
		case OP_MOD:
			mod_by, number := vm.popDataStack(), vm.popDataStack()
			result := modNumbers(number, mod_by)
			vm.pushDataStack(result)
		case OP_AND:
			and_with, number := vm.popDataStack(), vm.popDataStack()
			result := andNumbers(number, and_with)
			vm.pushDataStack(result)
		case OP_CALL:
			vm.pushCallStack(vm.Ip)
			vm.Ip = arg - 1
		case OP_RETURN:
			if len(vm.callStack) == 0 {
				return
			}
			vm.Ip = vm.popCallStack()
		case OP_PUSH:
			vm.pushDataStack(vm.Heap[arg])
		case OP_DUP:
			vm.pushDataStack(vm.dataStack[len(vm.dataStack) - int(arg) - 1])
		case OP_DROP:
			vm.dataStack = vm.dataStack[:len(vm.dataStack) - int(arg)]
		case OP_JUMP:
			vm.Ip = arg - 1
		case OP_JUMP_IF_NOT:
			value := vm.popDataStack()
			if value.DataType() == TYPE_INTEGER && value.(IntegerDatum).Int == 0 {
				vm.Ip = arg - 1
			}
		case OP_STORE:
			varName := vm.Heap[arg].(StringDatum).Str
			vm.variables[varName] = vm.popDataStack()
		case OP_FETCH:
			varName := vm.Heap[arg].(StringDatum).Str
			vm.pushDataStack(vm.variables[varName])
		}

		vm.Ip++
	}
}

func (vm *VirtualMachine) pushDataStack(datum Datum) {
	vm.dataStack = append(vm.dataStack, datum)
}

func (vm *VirtualMachine) popDataStack() Datum {
	datum := vm.dataStack[len(vm.dataStack) - 1]
  vm.dataStack = vm.dataStack[:len(vm.dataStack) - 1]
	return datum
}

func (vm *VirtualMachine) pushCallStack(address uint32) {
	vm.callStack = append(vm.callStack, address)
}

func (vm *VirtualMachine) popCallStack() uint32 {
	address := vm.callStack[len(vm.callStack) - 1]
  vm.callStack = vm.callStack[:len(vm.callStack) - 1]
	return address
}

func (vm *VirtualMachine) printDisassembly() {
	fmt.Println("Code disassembly:")
	for i, instruction := range vm.Code {
		opcode := uint8(instruction & 0xFF)
		arg := uint32(instruction >> 8)

		if uint32(i) == vm.Ip {
			fmt.Print("  IP> ")
		} else {
			fmt.Print("      ")
		}
		fmt.Printf("%04x: %08x   | %12s ", i, instruction, OpNames[opcode])

		switch opcode {
		case OP_PUSH:
		  printDatum(vm.Heap[arg], true)
		case OP_CALL:
			target := "<unknown routine>"
			for wordName, offset := range vm.Dict {
				if arg == offset {
					target = wordName
				}
			}
		  fmt.Printf("%s @ 0x%02x", target, arg)
		case OP_JUMP, OP_JUMP_IF_NOT:
		  fmt.Printf("%04x", arg)
		case OP_DUP, OP_DROP:
		  fmt.Print(arg)
		}

		for wordName, offset := range vm.Dict {
			if uint32(i) == offset {
				fmt.Printf("   [%s]", wordName)
			}
		}
		fmt.Println("")
	}
}

func printDatum(datum Datum, escaped bool) {
	switch datum.DataType() {
	case TYPE_INTEGER:
		fmt.Printf("%d", datum.(IntegerDatum).Int)
	case TYPE_STRING:
		if escaped {
			fmt.Printf("%#v", datum.(StringDatum).Str)
		} else {
			fmt.Printf("%s", datum.(StringDatum).Str)
		}
	default:
		panic(fmt.Sprintf("Can't print datum: %v", datum))
	}
}

func addNumbers(num1 Datum, num2 Datum) IntegerDatum {
	if num1.DataType() != TYPE_INTEGER || num2.DataType() != TYPE_INTEGER {
		panic("Can't add non-integer values!")
	}
	return IntegerDatum{num1.(IntegerDatum).Int + num2.(IntegerDatum).Int}
}

func modNumbers(num1 Datum, num2 Datum) IntegerDatum {
	if num1.DataType() != TYPE_INTEGER || num2.DataType() != TYPE_INTEGER {
		panic("Can't mod non-integer values!")
	}
	return IntegerDatum{num1.(IntegerDatum).Int % num2.(IntegerDatum).Int}
}

func andNumbers(num1 Datum, num2 Datum) IntegerDatum {
	if num1.DataType() != TYPE_INTEGER || num2.DataType() != TYPE_INTEGER {
		panic("Can't mod non-integer values!")
	}
	return IntegerDatum{num1.(IntegerDatum).Int & num2.(IntegerDatum).Int}
}
