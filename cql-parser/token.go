package cql

// Token represents a lexical token.
type Token int

// Tokens
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	COMMA              // ,
	Dot                // .
	LeftRoundBrackets  // (
	RightRoundBrackets // )
	LT                 // <
	GT                 // >

	// Keywords
	CreateTable
	Static
	PrimaryKey
)
