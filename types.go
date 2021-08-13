package main

const (
	OP_INVALID uint8 = iota
	OP_RETURN
	OP_PUSH
	OP_CALL
	OP_JUMP
	OP_JUMP_IF_NOT
	OP_PRINT
	OP_ADD
	OP_MOD
	OP_DUP
	OP_DROP
)

const (
  TYPE_VOID uint8 = iota
	TYPE_INTEGER
	TYPE_STRING
)

const (
	INTEGER_TOKEN uint8 = iota
	STRING_TOKEN
	KEYWORD_TOKEN
	FUNCALL_TOKEN
	EOF_TOKEN
)

// Too simple to be worth using an interface for.
type Token struct {
  TokenType uint8
	Int int64
	Str string
}

type Datum interface {
	DataType() uint8
}

type IntegerDatum struct {
	Int int64
}

type VoidDatum struct {}

type StringDatum struct {
	Str string
}

func (i VoidDatum) DataType() uint8 {
	return TYPE_VOID
}

func (i IntegerDatum) DataType() uint8 {
	return TYPE_INTEGER
}

func (i StringDatum) DataType() uint8 {
	return TYPE_STRING
}

type AbstractOp struct {
	Opcode uint8
	Arg uint32
	Datum Datum
}

type PackedOp uint32

type Word struct {
	Name string
	Ops []AbstractOp
}
