package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	cql "github.com/johnnywidth/cql-calculator/cql-parser"
	yaml "gopkg.in/yaml.v2"
)

func generateFromCQL(cqlString string) {
	b := bytes.NewReader([]byte(cqlString))
	p := cql.NewParser(b)
	r, err := p.Parse()
	if err != nil {
		panic(err)
	}

	m := Metadata{Name: r.TableName}

	fmt.Print("Enter rows count: ")
	var i int
	_, err = fmt.Scanf("%d", &i)
	if err != nil {
		panic(err)
	}

	m.Rows = i

	m.Partition = buildKeys(r.Colums, r.PK)
	m.Cluster = buildKeys(r.Colums, r.CK)
	m.Static = buildKeys(r.Colums, r.SK)

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

	data, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("generated.yaml", data, 0755)
	if err != nil {
		panic(err)
	}
}

func buildKeys(ac map[string]cql.Column, c map[string]cql.Column) []Column {
	nc := []Column{}
	for _, v := range c {
		t, ok := ac[v.Name]
		if !ok {
			log.Fatalf("Miss key in column %s", v.Name)
		}

		s := GetSizeByType(t.Name, t.Type)

		nc = append(nc, Column{Name: v.Name, Type: t.Type, Size: s})
	}

	return nc
}
