package main

import (
	"bufio"
	"io"
	"strconv"
)

type Parser struct {
	scanner *bufio.Scanner
}

// FIXME: Actual string parser that handles spaces and escaped characters in strings.
// FIXME: Words are currently case-sensitive, but should not be.
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
		for token := p.NextToken(); token.TokenType != KEYWORD_TOKEN || token.Str != ")"; token = p.NextToken() {
			if token.TokenType == EOF_TOKEN {
				panic("No matching ')' for '('!")
			}
		}
		return p.NextToken()
	default:
		return Token{FUNCALL_TOKEN, 0, s}
	}
}
