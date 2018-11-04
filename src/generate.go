package src

import (
	"bytes"
	"errors"

	cql "github.com/johnnywidth/cql-calculator/src/cql-parser"
)

func GenerateFromCQL(cqlString string, rn int) (Metadata, error) {
	b := bytes.NewReader([]byte(cqlString))
	p := cql.NewParser(b)
	r, err := p.Parse()
	if err != nil {
		return Metadata{}, err
	}

	if len(r.PK) == 0 {
		return Metadata{}, errors.New("There no partition keys")
	}

	m := Metadata{
		Rows: rn,
		Name: r.TableName,
	}

	m.Partition = populateColumnSize(r.PK)
	m.Cluster = populateColumnSize(r.CK)
	m.Static = populateColumnSize(r.SK)

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

		s := GetSizeByType(v.Name, v.Type)

		m.Column = append(m.Column, Column{Name: v.Name, Type: v.Type, Size: s})
	}

	return m, nil
}

func populateColumnSize(c map[string]cql.Column) []Column {
	nc := make([]Column, len(c))
	i := 0
	for _, v := range c {
		s := GetSizeByType(v.Name, v.Type)

		nc[i] = Column{Name: v.Name, Type: v.Type, Size: s}
		i++
	}

	return nc
}
