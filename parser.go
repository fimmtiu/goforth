package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	INTEGER_TOKEN uint8 = iota
	KEYWORD_TOKEN
	FUNCALL_TOKEN
	EOF_TOKEN
)

type Token struct {
  TokenType uint8
	Int int64
	Str string
}

type Parser struct {
	scanner *bufio.Scanner
}

func NewParser(data io.Reader) *Parser {
	p := Parser{bufio.NewScanner(data)}
	p.scanner.Split(bufio.ScanWords)
	return &p
}

// FIXME: Add an 'err' parameter to this instead of panicking.
func (p *Parser) NextToken() Token {
	if !p.scanner.Scan() {
		return Token{EOF_TOKEN, 0, ""}
	}
	s := p.scanner.Text()
	fmt.Printf("Read token: %v\n", s)

	if value, err := strconv.ParseInt(s, 10, 64); err == nil {
		return Token{INTEGER_TOKEN, value, ""}
	}

	switch s {
	case ":", ";", ")", "if", "then", "else":
		return Token{KEYWORD_TOKEN, 0, s}
	case "(":
		for token := p.NextToken(); token.TokenType != KEYWORD_TOKEN || token.Str != ")"; token = p.NextToken() {
			fmt.Printf("Skipping token: %v\n", token)
			if token.TokenType == EOF_TOKEN {
				panic("No matching ')' for '('!")
			}
		}
		return p.NextToken()
	default:
		return Token{FUNCALL_TOKEN, 0, s}
	}
}
