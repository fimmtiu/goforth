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
			t.Errorf("Expected token %d to be %v, but it was %v", i, tokens[i], token)
		}
		if token.TokenType == EOF_TOKEN {
			if i < len(tokens) - 1 {
				t.Errorf("Not enough tokens! Next token should have been %v", tokens[i + 1])
			}
			break
		}
	}
}

func TestEmptyInput(t *testing.T) {
	compareTokens(t, "", Token{EOF_TOKEN, 0, ""})
}

func TestIntegers(t *testing.T) {
	compareTokens(t, "1 31337 -7", Token{INTEGER_TOKEN, 1, ""}, Token{INTEGER_TOKEN, 31337, ""}, Token{INTEGER_TOKEN, -7, ""}, Token{EOF_TOKEN, 0, ""})
}

func TestStrings(t *testing.T) {
	compareTokens(t, `"1" "" "foo"`, Token{STRING_TOKEN, 0, "1"}, Token{STRING_TOKEN, 0, ""}, Token{STRING_TOKEN, 0, "foo"}, Token{EOF_TOKEN, 0, ""})
}

func TestIdentifiers(t *testing.T) {
	compareTokens(t, "a A foo? ?bar - ", Token{FUNCALL_TOKEN, 0, "a"}, Token{FUNCALL_TOKEN, 0, "A"}, Token{FUNCALL_TOKEN, 0, "foo?"}, Token{FUNCALL_TOKEN, 0, "?bar"}, Token{FUNCALL_TOKEN, 0, "-"}, Token{EOF_TOKEN, 0, ""})
}

func TestComments(t *testing.T) {
	compareTokens(t, "2 ( I like pie ) .", Token{INTEGER_TOKEN, 2, ""}, Token{FUNCALL_TOKEN, 0, "."}, Token{EOF_TOKEN, 0, ""})
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
		if token.TokenType == EOF_TOKEN {
			break
		}
	}
}
