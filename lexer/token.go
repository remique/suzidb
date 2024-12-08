package lexer

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"

	// Keywords
	SELECT = "SELECT"
	INSERT = "INSERT"
	WHERE  = "WHERE"
	FROM   = "FROM"
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
