package cql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	COMMA              // ,
	LeftRoundBrackets  // (
	RightRoundBrackets // )

	// Keywords
	CREATE_TABLE
	STATIC
	PRIMARY_KEY
)
