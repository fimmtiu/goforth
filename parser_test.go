package main

import (
	"fmt"
	"strings"
	"testing"
)

func compareTokens(t *testing.T, code string, tokens ...Token) {
	parser := NewParser(strings.NewReader(code))

	for i := 0 ;; i++ {
		token := parser.NextToken()
		fmt.Println(token)
		if i >= len(tokens) {
			t.Errorf("Too many tokens! Extra token was %v", token)
		}
		if token != tokens[i] {
			t.Errorf("Unexpected token: %v", token)
		}
		if token.TokenType == EOF {
			if i < len(tokens) - 1 {
				t.Errorf("Not enough tokens! Next token should have been %v", tokens[i + 1])
			}
			break
		}
	}
}

func TestEmptyInput(t *testing.T) {
	compareTokens(t, "", Token{EOF, 0, ""})
}

func TestIntegers(t *testing.T) {
	compareTokens(t, "1 31337 -7", Token{INTEGER, 1, ""}, Token{INTEGER, 31337, ""}, Token{INTEGER, -7, ""}, Token{EOF, 0, ""})
}

func TestIdentifiers(t *testing.T) {
	compareTokens(t, "a A foo? ?bar - ", Token{FUNCALL, 0, "a"}, Token{FUNCALL, 0, "A"}, Token{FUNCALL, 0, "foo?"}, Token{FUNCALL, 0, "?bar"}, Token{FUNCALL, 0, "-"}, Token{EOF, 0, ""})
}

func TestComments(t *testing.T) {
	compareTokens(t, "2 ( I like pie ) .", Token{INTEGER, 2, ""}, Token{FUNCALL, 0, "."}, Token{EOF, 0, ""})
}

func TestUnboundedComment(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
				t.Errorf("The code did not panic")
		}
	}()

	parser := NewParser(strings.NewReader("1 ( 2")) // Should panic with "No matching ')'" error
	for {
		token := parser.NextToken()
		if token.TokenType == EOF {
			break
		}
	}
}
