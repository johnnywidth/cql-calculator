package cql

import (
	"fmt"
	"io"
)

// Statement represents a CQL CREATE TABLE statement.
type Statement struct {
	TableName string
	Colums    map[string]Column
	PK        map[string]Column
	CK        map[string]Column
	SK        map[string]Column
}

// Column column data
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

// Parse parses a CQL CREATE TABLE statement.
func (p *Parser) Parse() (*Statement, error) {
	stmt := &Statement{
		Colums: make(map[string]Column),
		PK:     make(map[string]Column),
		CK:     make(map[string]Column),
		SK:     make(map[string]Column),
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CreateTable {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != CreateTable {
		return nil, fmt.Errorf("found %q, expected TABLE", lit)
	}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	tok, lit = p.scanIgnoreWhitespace()
	if tok == Dot {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected table name", lit)
		}
		stmt.TableName = lit

		tok, lit = p.scanIgnoreWhitespace()
	}

	if tok != LeftRoundBrackets {
		return nil, fmt.Errorf("found %q, expected lrb", lit)
	}

	err := p.regularColumns(stmt)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) regularColumns(stmt *Statement) error {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == PrimaryKey {
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
		if tok3 == Static {
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

func (p *Parser) pkColumns(stmt *Statement) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != PrimaryKey {
		return fmt.Errorf("PK. found %q, expected primary key", lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok != LeftRoundBrackets {
		return fmt.Errorf("PK. found %q, expected left round brackets", lit)
	}

	tok, _ = p.scanIgnoreWhitespace()
	if tok == LeftRoundBrackets {
		return p.pkCompositeColumns(stmt)
	}

	p.unscan()

	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("PK. found %q, expected column name", lit)
		}

		stmt.PK[lit] = Column{Name: lit, Type: stmt.Colums[lit].Type}
		break
	}

	err := p.ckColumns(stmt)

	return err
}

func (p *Parser) pkCompositeColumns(stmt *Statement) error {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("Composite PK. found %q, expected column name", lit)
		}

		stmt.PK[lit] = Column{Name: lit, Type: stmt.Colums[lit].Type}

		tok, _ = p.scanIgnoreWhitespace()
		if tok == RightRoundBrackets {
			break
		}
		if tok != COMMA {
			p.unscan()
			break
		}
	}

	err := p.ckColumns(stmt)

	return err
}

func (p *Parser) ckColumns(stmt *Statement) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok == RightRoundBrackets {
		return nil
	}
	if tok != COMMA {
		return fmt.Errorf("Cluster. found %q, expected rrb or comma", lit)
	}

	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("Cluster. found %q, expected column name", lit)
		}

		stmt.CK[lit] = Column{Name: lit, Type: stmt.Colums[lit].Type}

		tok, _ = p.scanIgnoreWhitespace()
		if tok != COMMA {
			p.unscan()
			break
		}
	}

	return nil
}
