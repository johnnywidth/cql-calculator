package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/johnnywidth/cql-calculator/src"
	cql "github.com/johnnywidth/cql-calculator/src/cql-parser"
	yaml "gopkg.in/yaml.v2"
)

func generateFromCQL(fileName, cqlString string) {
	b := bytes.NewReader([]byte(cqlString))
	p := cql.NewParser(b)
	r, err := p.Parse()
	if err != nil {
		panic(err)
	}

	if len(r.PK) == 0 {
		panic("There no partition keys")
	}

	m := src.Metadata{Name: r.TableName}

	fmt.Print("Enter rows count per one partition: ")
	var i int
	_, err = fmt.Scanf("%d", &i)
	if err != nil {
		panic(err)
	}

	m.Rows = i

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

		s := src.GetSizeByType(v.Name, v.Type)

		m.Column = append(m.Column, src.Column{Name: v.Name, Type: v.Type, Size: s})
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, data, 0755)
	if err != nil {
		panic(err)
	}
}

func populateColumnSize(c map[string]cql.Column) []src.Column {
	nc := make([]src.Column, len(c))
	i := 0
	for _, v := range c {
		s := src.GetSizeByType(v.Name, v.Type)

		nc[i] = src.Column{Name: v.Name, Type: v.Type, Size: s}
		i++
	}

	return nc
}
