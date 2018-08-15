package main

import (
	"bytes"
	"io/ioutil"

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

	for _, v := range r.PK {
		m.Partition = append(m.Partition, Column{Name: v.Name, Type: v.Type})
	}

	for _, v := range r.CK {
		m.Cluster = append(m.Cluster, Column{Name: v.Name, Type: v.Type})
	}

	for _, v := range r.SK {
		m.Static = append(m.Static, Column{Name: v.Name, Type: v.Type})
	}

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
		m.Column = append(m.Column, Column{Name: v.Name, Type: v.Type})
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
