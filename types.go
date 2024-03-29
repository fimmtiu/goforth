package main

const (
	OP_INVALID uint8 = iota   // 00
	OP_RETURN                 // 01
	OP_PUSH                   // 02
	OP_CALL                   // 03
	OP_JUMP                   // 04
	OP_JUMP_IF_NOT            // 05
	OP_PRINT                  // 06
	OP_ADD                    // 07
	OP_MOD                    // 08
	OP_DUP                    // 09
	OP_DROP                   // 0a
	OP_AND                    // 0b
	OP_STORE                  // 0c
	OP_FETCH                  // 0d
)

var OpNames = []string{
	"INVALID",
	"RETURN",
	"PUSH",
	"CALL",
	"JUMP",
	"JUMP_IF_NOT",
	"PRINT",
	"ADD",
	"MOD",
	"DUP",
	"DROP",
	"AND",
	"STORE",
	"FETCH",
}

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
