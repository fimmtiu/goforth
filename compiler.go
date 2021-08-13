package main

import (
	"fmt"
	"io"
)

type Compiler struct {
	parser *Parser
	vm *VirtualMachine
  compiling bool
	words []Word
	pushedBackToken *Token
}

func (w *Word) Finish() {
	w.Ops = append(w.Ops, AbstractOp{OP_RETURN, 0, VoidDatum{}})
}

func NewCompiler(vm *VirtualMachine) *Compiler {
	return &Compiler{nil, vm, false, []Word{{"top-level code", []AbstractOp{}}}, nil}
}

func (c *Compiler) ReadToken() Token {
  if c.pushedBackToken != nil {
		token := c.pushedBackToken
		c.pushedBackToken = nil
		return *token
	}
	return c.parser.NextToken()
}

func (c *Compiler) UnreadToken(t Token) {
	if c.pushedBackToken != nil {
		panic(fmt.Sprintf("WTF: Token %v already pushed, but tried to push %v", *c.pushedBackToken, t))
	}
	c.pushedBackToken = &t
}

// The first byte of the uint32 is the opcode; the remaining 3 bytes are some sort of argument to the instruction.
func (c *Compiler) convertToPackedOp(word Word, op AbstractOp, opIndex int) PackedOp {
	var offset uint32 = 0

	switch op.Opcode {
	case OP_CALL:
		word_name := op.Datum.(StringDatum).Str
		offset = c.vm.Dict[word_name]

	case OP_PUSH:
		c.vm.Heap = append(c.vm.Heap, op.Datum)
		offset = uint32(len(c.vm.Heap)) - 1

	case OP_JUMP, OP_JUMP_IF_NOT:
		offset = c.vm.Dict[word.Name] + uint32(opIndex) + uint32(op.Arg)
	}
	return PackedOp(uint32(op.Opcode) | (offset << 8))
}

// FIXME: Actual error handling
// FIXME: I don't like that Compiler reaches into VM like this.
func (c *Compiler) LoadCode(code io.Reader) {
	c.parser = NewParser(code)
	c.words[0].Ops = c.Compile()
	c.words[0].Finish()

	// Populate the dictionary with the starting offsets of each word in the code array.
	var offset uint32 = uint32(len(c.vm.Code))
	for _, word := range c.words {
		c.vm.Dict[word.Name] = offset
		offset += uint32(len(word.Ops))
	}

	// Convert all the AbstractOps to PackedOps and store them in the VM.
	for _, word := range c.words {
		packedOps := []PackedOp{}

		for i, op := range word.Ops {
			packedOps = append(packedOps, c.convertToPackedOp(word, op, i))
		}
		c.vm.Code = append(c.vm.Code, packedOps...)
	}

	c.parser = nil
}

// FIXME: This function is long. Break it up?
// FIXME: Actual error handling instead of panics.
func (c *Compiler) Compile(stopwords ...string) []AbstractOp {
	ops := []AbstractOp{}

	for {
		token := c.ReadToken()

		switch token.TokenType {
		case KEYWORD_TOKEN:
			for _, stopword := range stopwords {
				if token.Str == stopword {
					c.UnreadToken(token)
					return ops
				}
			}

			switch token.Str {
			case ":":
				c.defineWord()
			case ";":
				panic("Can't use ';' outside of a word definition!")
			case "if":
				ops = append(ops, c.compileIf()...)
			case "else", "then":
				panic(fmt.Sprintf("Can't have '%s' without a matching 'if'!", token.Str))
			default:
				panic(fmt.Sprintf("Unknown keyword: %v", token))
			}

		case INTEGER_TOKEN:
			ops = append(ops, AbstractOp{OP_PUSH, 0, IntegerDatum{token.Int}})

		case STRING_TOKEN:
			ops = append(ops, AbstractOp{OP_PUSH, 0, StringDatum{token.Str}})

		case FUNCALL_TOKEN:
			switch token.Str {
			case ".":
				ops = append(ops, AbstractOp{OP_PRINT, 0, VoidDatum{}})
			case "+":
				ops = append(ops, AbstractOp{OP_ADD, 0, VoidDatum{}})
			default:
				ops = append(ops, AbstractOp{OP_CALL, 0, StringDatum{token.Str}})
			}

		case EOF_TOKEN:
			return ops

		default:
			panic(fmt.Sprintf("Unknown token type: %v", token))
		}
	}
}

func (c *Compiler) defineWord() {
	if c.compiling {
		panic("Can't nest word definitions!")
	}
	c.compiling = true

	nameToken := c.ReadToken()
	if nameToken.TokenType != FUNCALL_TOKEN {
		panic(fmt.Sprintf("'%v' isn't a valid word name!", nameToken))
	}

	word := Word{nameToken.Str, c.Compile(";")}

	// Consume the trailing ';' token
	terminator := c.ReadToken()
	if terminator.TokenType != KEYWORD_TOKEN || terminator.Str != ";" {
		panic(fmt.Sprintf("EOF during word definition for '%v'!", nameToken.Str))
	}

	word.Finish()
	c.words = append(c.words, word)
	c.compiling = false
}

func (c *Compiler) compileIf() []AbstractOp {
	ops := []AbstractOp{}
	true_branch := c.Compile("else", "then")
	false_branch := []AbstractOp{}

	nextToken := c.ReadToken()
	if nextToken.TokenType == KEYWORD_TOKEN && nextToken.Str == "else" {
		false_branch = c.Compile("then")
		nextToken = c.ReadToken()
	}

	if nextToken.TokenType != KEYWORD_TOKEN || nextToken.Str != "then" {
		panic("Improperly terminated 'if' statement!")
	}

	if len(false_branch) > 0 {
		ops = append(ops, AbstractOp{OP_JUMP_IF_NOT, uint32(len(true_branch) + 2), VoidDatum{}})
		ops = append(ops, true_branch...)
		ops = append(ops, AbstractOp{OP_JUMP, uint32(len(false_branch) + 1), VoidDatum{}})
		ops = append(ops, false_branch...)
	} else {
		ops = append(ops, AbstractOp{OP_JUMP_IF_NOT, uint32(len(true_branch) + 1), VoidDatum{}})
		ops = append(ops, true_branch...)
	}
	return ops
}
