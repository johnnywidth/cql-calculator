package src

import (
	"bytes"
	"errors"

	cql "github.com/johnnywidth/cql-calculator/src/cql-parser"
)

// Metadata table metadata
type Metadata struct {
	Name      string   `yaml:"name"`
	Rows      int      `yaml:"rows"`
	Partition []Column `yaml:"partition"`
	Cluster   []Column `yaml:"cluster"`
	Static    []Column `yaml:"static"`
	Column    []Column `yaml:"column"`

	customSizes map[string]CustomSize
}

// Column table column data
type Column struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Size int    `yaml:"size"`
}

// PopulateTableMetadata populate table metadata from cql query
func PopulateTableMetadata(cqlString string, rn int) (Metadata, error) {
	b := bytes.NewReader([]byte(cqlString))
	p := cql.NewParser(b)
	r, err := p.Parse()
	if err != nil {
		return Metadata{}, err
	}

	if len(r.PK) == 0 {
		return Metadata{}, errors.New("missed partition key(s)")
	}

	m := Metadata{
		Rows: rn,
		Name: r.TableName,
	}

	m.Partition = m.populateColumnSize(r.PK)
	m.Cluster = m.populateColumnSize(r.CK)
	m.Static = m.populateColumnSize(r.SK)
	m.populateRegularColumnSize(r)

	return m, nil
}

func (m *Metadata) populateColumnSize(c map[string]cql.Column) []Column {
	nc := make([]Column, 0)
	for _, v := range c {
		s, err := GetSizeByType(v.Name, v.Type)
		if err != nil {
			m.addCustomSize(v)
		}

		nc = append(nc, Column{Name: v.Name, Type: v.Type, Size: s})
	}

	return nc
}

func (m *Metadata) populateRegularColumnSize(r *cql.Statement) {
	for _, v := range r.Colums {
		if _, ok := r.PK[v.Name]; ok {
			continue
		}
		if _, ok := r.CK[v.Name]; ok {
			continue
		}
		if _, ok := r.SK[v.Name]; ok {
			continue
		}

		s, err := GetSizeByType(v.Name, v.Type)
		if err != nil {
			m.addCustomSize(v)
		}

		m.Column = append(m.Column, Column{Name: v.Name, Type: v.Type, Size: s})
	}
}
