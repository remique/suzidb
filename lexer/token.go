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
	DOT       = "."

	// Keywords
	SELECT  = "SELECT"
	INSERT  = "INSERT"
	WHERE   = "WHERE"
	FROM    = "FROM"
	INTO    = "INTO"
	VALUES  = "VALUES"
	CREATE  = "CREATE"
	TABLE   = "TABLE"
	PRIMARY = "PRIMARY"
	KEY     = "KEY"
	NOT     = "NOT"
	NULL    = "NULL"

	// Join Keywords
	JOIN  = "JOIN"
	ON    = "ON"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	INNER = "INNER"

	// (Keywords) To denote a type in CREATE TABLE statement
	TEXT_TYPE = "TEXT_TYPE"
	INT_TYPE  = "INT_TYPE"

	L_PAREN = "("
	R_PAREN = ")"

	// Illegal & EOF
	ILLEGAL
	EOF
)

// A map of all reserved keywords
var keywords = map[string]TokenType{
	"select":  SELECT,
	"insert":  INSERT,
	"from":    FROM,
	"where":   WHERE,
	"into":    INTO,
	"values":  VALUES,
	"create":  CREATE,
	"table":   TABLE,
	"primary": PRIMARY,
	"key":     KEY,
	"text":    TEXT_TYPE,
	"int":     INT_TYPE,
	"not":     NOT,
	"null":    NULL,
	"join":    JOIN,
	"left":    LEFT,
	"right":   RIGHT,
	"inner":   INNER,
	"on":      ON,
}

type Token struct {
	TokenType TokenType
	Literal   string
}

func NewToken(tokenType TokenType, literal string) Token {
	return Token{
		TokenType: tokenType,
		Literal:   literal,
	}
}
