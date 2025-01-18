package lexer

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"

	// For values assigned within tables, eg. 'Alice', 'Bob'
	STRING = "STRING"
	INT    = "INT"

	EQUALS    = "="
	PLUS      = "+"
	STAR      = "*"
	SEMICOLON = ";"
	COMMA     = ","

	// Keywords
	SELECT = "SELECT"
	INSERT = "INSERT"
	WHERE  = "WHERE"
	FROM   = "FROM"
	INTO   = "INTO"
	VALUES = "VALUES"

	L_PAREN = "("
	R_PAREN = ")"

	// Illegal & EOF
	ILLEGAL
	EOF
)

var keywords = map[string]TokenType{
	"select": SELECT,
	"insert": INSERT,
	"from":   FROM,
	"where":  WHERE,
	"into":   INTO,
	"values": VALUES,
}

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
