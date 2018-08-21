package cql

import (
	"fmt"
	"io"
)

// SelectStatement represents a SQL CREATE TABLE statement.
type SelectStatement struct {
	TableName string
	Colums    map[string]Column
	PK        map[string]Column
	CK        map[string]Column
	SK        map[string]Column
}

type Column struct {
	Name string
	Type string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a SQL CREATE TABLE statement.
func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{
		Colums: make(map[string]Column),
		PK:     make(map[string]Column),
		CK:     make(map[string]Column),
		SK:     make(map[string]Column),
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CREATE_TABLE {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CREATE_TABLE {
		return nil, fmt.Errorf("found %q, expected TABLE", lit)
	}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	if tok, lit := p.scanIgnoreWhitespace(); tok != LeftRoundBrackets {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}

	err := p.regularColumns(stmt)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) regularColumns(stmt *SelectStatement) error {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == PRIMARY_KEY {
			err := p.pkColumns(stmt)
			if err != nil {
				return err
			}
			break
		}

		if tok != IDENT {
			return fmt.Errorf("found %q, expected column name", lit)
		}

		tok2, lit2 := p.scanIgnoreWhitespace()
		if tok2 != IDENT {
			return fmt.Errorf("found %q, expected type", lit2)
		}

		stmt.Colums[lit] = Column{Name: lit, Type: lit2}

		tok3, _ := p.scanIgnoreWhitespace()
		if tok3 == STATIC {
			stmt.SK[lit] = Column{Name: lit, Type: lit2}
		} else {
			p.unscan()
		}

		tok, _ = p.scanIgnoreWhitespace()
		if tok != COMMA {
			p.unscan()
			break
		}
	}

	return nil
}

func (p *Parser) pkColumns(stmt *SelectStatement) error {
	c := make(map[string]Column)

	tok, lit := p.scanIgnoreWhitespace()
	if tok != PRIMARY_KEY {
		return fmt.Errorf("found %q, expected primary key", lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok != LeftRoundBrackets {
		return fmt.Errorf("found %q, expected left round brackets", lit)
	}

	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("found %q, expected column name", lit)
		}

		c[lit] = Column{Name: lit}

		// TODO: if PKs more than one in `()`
		break

		tok, _ = p.scanIgnoreWhitespace()
		if tok != COMMA {
			p.unscan()
			break
		}
	}

	stmt.PK = c

	err := p.ckColumns(stmt)

	return err
}

func (p *Parser) ckColumns(stmt *SelectStatement) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != RightRoundBrackets && tok != COMMA {
		return fmt.Errorf("found %q, expected rrb or comma", lit)
	}

	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("found %q, expected column name", lit)
		}

		stmt.CK[lit] = Column{Name: lit}

		tok, _ = p.scanIgnoreWhitespace()
		if tok != COMMA {
			p.unscan()
			break
		}
	}

	return nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
