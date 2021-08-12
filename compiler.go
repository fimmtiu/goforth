package main
type Op struct {
	Opcode uint8
	Arg uint32
}
type Word struct {
	Name string
	Ops []Op
}

