package main

import "fmt"

const (
	OP_RETURN uint8 = iota
	OP_PUSH
	OP_CALL
	OP_JUMP
	OP_JUMP_IF_NOT
	OP_PRINT
)

const (
  TYPE_VOID = iota
	TYPE_INTEGER
	TYPE_STRING
)

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

type Op struct {
	Opcode uint8
	Arg uint32
	Datum Datum
}

type Word struct {
	Name string
	Ops []Op
}

type Compiler struct {
	parser *Parser
	// vm VirtualMachine
  compiling bool
	words []Word
	pushedBackToken *Token
}

func (w *Word) Finish() {
	w.Ops = append(w.Ops, Op{OP_RETURN, 0, VoidDatum{}})
}

func NewCompiler(p *Parser) *Compiler {
	return &Compiler{
		p, false, []Word{}, nil,
	}
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

// FIXME: This function is long. Break it up?
// FIXME: Actual error handling instead of panics.
func (c *Compiler) Compile(stopwords ...string) []Op {
	ops := []Op{}

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
				c.compileWord()
			case ";":
				panic("Can't use ';' outside of a word definition!")
			case "if":
				// FIXME not implemented yet
			default:
				panic(fmt.Sprintf("Unknown keyword: %v", token))
			}

		case INTEGER_TOKEN:
			ops = append(ops, Op{OP_PUSH, 0, IntegerDatum{token.Int}})

		case FUNCALL_TOKEN:
			switch token.Str {
			case ".":
				ops = append(ops, Op{OP_PRINT, 0, VoidDatum{}})
			default:
				ops = append(ops, Op{OP_CALL, 0, StringDatum{token.Str}})
			}

		case EOF_TOKEN:
			return ops

		default:
			panic(fmt.Sprintf("Unknown token type: %v", token))
		}
	}
}

func (c *Compiler) compileWord() {
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
