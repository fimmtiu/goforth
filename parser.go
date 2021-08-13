package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	scanner *bufio.Scanner
	pushedBackToken *Token
}

// FIXME: Actual string parser that handles spaces and escaped characters in strings.
// FIXME: Words are currently case-sensitive, but should not be.
func NewParser(data io.Reader) *Parser {
	p := Parser{bufio.NewScanner(data), nil}
	p.scanner.Split(bufio.ScanWords)
	return &p
}

func (p *Parser) ReadToken() Token {
  if p.pushedBackToken != nil {
		token := p.pushedBackToken
		p.pushedBackToken = nil
		return *token
	}
	return p.nextToken()
}

func (p *Parser) UnreadToken(t Token) {
	if p.pushedBackToken != nil {
		panic(fmt.Sprintf("WTF: Token %v already pushed, but tried to push %v", *p.pushedBackToken, t))
	}
	p.pushedBackToken = &t
}

func (p *Parser) PeekToken() Token {
	token := p.ReadToken()
	p.UnreadToken(token)
	return token
}

// FIXME: Add an 'err' parameter to this instead of panicking.
func (p *Parser) nextToken() Token {
	if !p.scanner.Scan() {
		return Token{EOF_TOKEN, 0, ""}
	}
	s := p.scanner.Text()

	if value, err := strconv.ParseInt(s, 10, 64); err == nil {
		return Token{INTEGER_TOKEN, value, ""}
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		if s == `"\n"` { // Someday I'll parse strings correctly!
			return Token{STRING_TOKEN, 0, "\n"}
		} else {
			return Token{STRING_TOKEN, 0, s[1:len(s)-1]}
		}
	}

	switch s {
	case ":", ";", ")", "if", "then", "else":
		return Token{KEYWORD_TOKEN, 0, s}
	case "(":
		for token := p.ReadToken(); token.TokenType != KEYWORD_TOKEN || token.Str != ")"; token = p.ReadToken() {
			if token.TokenType == EOF_TOKEN {
				panic("No matching ')' for '('!")
			}
		}
		return p.ReadToken()
	default:
		return Token{FUNCALL_TOKEN, 0, s}
	}
}
