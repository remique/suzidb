package lexer

type Lexer struct {
	input           string
	currentPosition int
	ch              byte
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:           input,
		currentPosition: 0,
	}
}

func (l *Lexer) readChar() {
	if l.currentPosition >= len(l.input) {
		l.ch = '0'
	} else {
		l.ch = l.input[l.currentPosition]
	}

	l.currentPosition += 1
}
