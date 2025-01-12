package tokens

type Token int

const (
	ERROR Token = iota
	EOF
	NEWLINE
	IDENTIFIER
	NUMBER

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	COMMA // ,
	COLON // :

	// Keywords
	COMPONENT
	TEST
	DEFINE
	SET
	ASSERT
	EXPECT
	IS
	FOR
	FROM
	TO
)
