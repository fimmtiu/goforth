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

type Op struct {
	Opcode uint8
	Arg uint32
	Token Token
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
	w.Ops = append(w.Ops, Op{OP_RETURN, 0, Token{}})
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
				if token.String == stopword {
					c.UnreadToken(token)
					return ops
				}
			}
			switch token.String {
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
			ops = append(ops, Op{OP_PUSH, 0, token})
		case FUNCALL_TOKEN:
			switch token.String {
			case ".":
				ops = append(ops, Op{OP_PRINT, 0, Token{}})
			default:
				ops = append(ops, Op{OP_CALL, 0, token})
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

	word := Word{nameToken.String, c.Compile(";")}

	// Consume the trailing ';' token
	terminator := c.ReadToken()
	if terminator.TokenType != KEYWORD_TOKEN || terminator.String != ";" {
		panic(fmt.Sprintf("EOF during word definition for '%v'!", nameToken.String))
	}

	word.Finish()
	c.words = append(c.words, word)
	c.compiling = false
}
