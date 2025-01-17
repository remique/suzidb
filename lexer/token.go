package lexer

type TokenType string

var keywords = map[string]TokenType{
	"select": SELECT,
	"insert": INSERT,
}

const (
	IDENTIFIER = "IDENTIFIER"

	EQUALS = "="
	PLUS   = "+"

	// Keywords
	SELECT = "SELECT"
	INSERT = "INSERT"
	WHERE  = "WHERE"
	FROM   = "FROM"

	// Illegal & EOF
	ILLEGAL
	EOF
)

type Token struct {
	tokenType TokenType
	literal   string
}

func newToken(tokenType TokenType, literal string) Token {
	return Token{
		tokenType: tokenType,
		literal:   literal,
	}
}
