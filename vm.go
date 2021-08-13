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
}

func NewVirtualMachine() *VirtualMachine {
	var vm VirtualMachine
	vm.Dict = make(map[string]uint32)
	return &vm
}

func (vm *VirtualMachine) Run() {
	// fmt.Println("Code:")
	// for i, op := range vm.Code {
	// 	fmt.Printf("    %02d: %08x\n", i, op)
	// }

	for {
		instruction := vm.Code[vm.Ip]
		opcode := uint8(instruction & 0xFF)
		arg := uint32(instruction >> 8)
		// fmt.Printf("Executing: [%d] opcode %d, arg %d\n", vm.Ip, opcode, arg)

		switch opcode {
		case OP_PRINT:
			printDatum(vm.popDataStack())
		case OP_ADD:
			result := addNumbers(vm.popDataStack(), vm.popDataStack())
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
		case OP_JUMP:
			vm.Ip = arg - 1
		case OP_JUMP_IF_NOT:
			value := vm.popDataStack()
			if value.DataType() == TYPE_INTEGER && value.(IntegerDatum).Int == 0 {
				vm.Ip = arg - 1
			}
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

func printDatum(datum Datum) {
	switch datum.DataType() {
	case TYPE_INTEGER:
		fmt.Printf("%d", datum.(IntegerDatum).Int)
	case TYPE_STRING:
		fmt.Printf("%s", datum.(StringDatum).Str)
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
