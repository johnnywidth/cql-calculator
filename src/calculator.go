package src

import "fmt"

type Metadata struct {
	Name      string   `yaml:"name"`
	Rows      int      `yaml:"rows"`
	Partition []Column `yaml:"partition"`
	Cluster   []Column `yaml:"cluster"`
	Static    []Column `yaml:"static"`
	Column    []Column `yaml:"column"`
}

type Column struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Size int    `yaml:"size"`
}

type NOV struct {
	Metadata Metadata

	result  int
	formula string
}

func (p *NOV) Calculate() {
	allColumns := len(p.Metadata.Column) + len(p.Metadata.Partition) + len(p.Metadata.Cluster) + len(p.Metadata.Static)
	primaryColumns := len(p.Metadata.Partition) + len(p.Metadata.Cluster)
	staticColumns := len(p.Metadata.Static)

	p.formula = fmt.Sprintf("(%d * (%d - %d - %d) + %d)", p.Metadata.Rows, allColumns, primaryColumns, staticColumns, staticColumns)

	p.result = p.Metadata.Rows*(allColumns-primaryColumns-staticColumns) + staticColumns
}

func (p *NOV) GetResult() int {
	return p.result
}

func (p *NOV) GetFormula() string {
	return p.formula
}

type PDS struct {
	Metadata Metadata
	NOV      int

	formula string
	result  int
}

func (p *PDS) Calculate() {
	var sofPK int
	for _, v := range p.Metadata.Partition {
		sofPK += v.Size
	}

	var sofS int
	for _, v := range p.Metadata.Static {
		sofS += v.Size
	}

	var sofCK int
	for _, v := range p.Metadata.Cluster {
		sofCK += v.Size
	}

	var ek int
	for _, v := range p.Metadata.Column {
		ek += v.Size + sofCK
	}

	p.formula = fmt.Sprintf("(%d + %d + (%d * %d) + (8 * %d))", sofPK, sofS, p.Metadata.Rows, ek, p.NOV)

	p.result = sofPK + sofS + (p.Metadata.Rows * ek) + (8 * p.NOV)
}

func (p *PDS) GetResult() int {
	return p.result
}

func (p *PDS) GetFormula() string {
	return p.formula
}
