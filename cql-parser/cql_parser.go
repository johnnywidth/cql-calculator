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

	err := p.createTable()
	if err != nil {
		return nil, err
	}

	// Parse table_name or database.table_name
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

	// Starts from `(` to find columns
	if tok != LeftRoundBrackets {
		return nil, fmt.Errorf("found %q, expected '('", lit)
	}

	err = p.regularColumns(stmt)

	return stmt, err
}

func (p *Parser) createTable() error {
	// Starts from "CREATE TABLE"
	if tok, lit := p.scanIgnoreWhitespace(); tok != CreateTable {
		return fmt.Errorf("found %q, expected `CREATE TABLE`", lit)
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != CreateTable {
		return fmt.Errorf("found %q, expected `CREATE TABLE`", lit)
	}

	// Consider optional IF NOT EXISTS after CREATE TABLE
	if tok, _ := p.scanIgnoreWhitespace(); tok != IfNotExists {
		p.unscan()
	} else {
		if tok, lit := p.scanIgnoreWhitespace(); tok != IfNotExists {
			return fmt.Errorf("found %q, expected 'IF NOT EXISTS'", lit)
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok != IfNotExists {
			return fmt.Errorf("found %q, expected 'IF NOT EXISTS'", lit)
		}
	}

	return nil
}

func (p *Parser) regularColumns(stmt *Statement) error {
	for {
		tok, columnName := p.scanIgnoreWhitespace()
		if tok == PrimaryKey {
			err := p.pkColumns(stmt)
			if err != nil {
				return err
			}
			break
		}

		if tok != IDENT {
			return fmt.Errorf("parse regular columns failed. found %q, expected column name", columnName)
		}

		stmt.Colums[columnName] = Column{Name: columnName, Type: p.columnType(stmt, columnName)}
	}

	return nil
}

func (p *Parser) columnType(stmt *Statement, name string) string {
	var cType string
	var collectionBracket int

	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == COMMA && collectionBracket == 0 {
			break
		}
		if tok == LT {
			collectionBracket++
		}
		if tok == GT {
			collectionBracket--
		}

		if tok == Static {
			stmt.SK[name] = Column{Name: name, Type: cType}
			continue
		}

		cType += lit
	}

	return cType
}

func (p *Parser) pkColumns(stmt *Statement) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != PrimaryKey {
		return fmt.Errorf("parse primary key. found %q, expected primary key", lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok != LeftRoundBrackets {
		return fmt.Errorf("parse primary key. found %q, expected '('", lit)
	}

	tok, _ = p.scanIgnoreWhitespace()
	if tok == LeftRoundBrackets {
		err := p.pkCompositeColumns(stmt)
		if err != nil {
			return err
		}
	} else {
		p.unscan()

		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("parse primary key. found %q, expected column name", lit)
		}

		stmt.PK[lit] = Column{Name: lit, Type: stmt.Colums[lit].Type}
	}

	return p.ckColumns(stmt)
}

func (p *Parser) pkCompositeColumns(stmt *Statement) error {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("parse composite partition key failed. found %q, expected column name", lit)
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

	return nil
}

func (p *Parser) ckColumns(stmt *Statement) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok == RightRoundBrackets {
		return nil
	}
	if tok != COMMA {
		return fmt.Errorf("parse cluster columns failed. found %q, expected ')' or ','", lit)
	}

	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return fmt.Errorf("parse cluster columns failed. found %q, expected column name", lit)
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
