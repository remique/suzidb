package lexer

import (
	"testing"
)

func TestReadChar(t *testing.T) {
	lexer := NewLexer("abc")

	tests := []struct {
		expectedChar byte
	}{
		{'a'},
		{'b'},
		{'c'},
		{'0'},
		{'0'},
	}

	for _, test := range tests {
		lexer.readChar()

		if lexer.ch != test.expectedChar {
			t.Fatalf("Expected %q, got %q", test.expectedChar, lexer.ch)
		}
	}
}
