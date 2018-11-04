package src

import "fmt"

// Metadata table metadata
type Metadata struct {
	Name      string   `yaml:"name"`
	Rows      int      `yaml:"rows"`
	Partition []Column `yaml:"partition"`
	Cluster   []Column `yaml:"cluster"`
	Static    []Column `yaml:"static"`
	Column    []Column `yaml:"column"`
}

// Column table column data
type Column struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Size int    `yaml:"size"`
}

// NOV number of values data
type NOV struct {
	Metadata Metadata

	result  int
	formula string
}

// Calculate calculate number of values
func (p *NOV) Calculate() {
	allColumns := len(p.Metadata.Column) + len(p.Metadata.Partition) + len(p.Metadata.Cluster) + len(p.Metadata.Static)
	primaryColumns := len(p.Metadata.Partition) + len(p.Metadata.Cluster)
	staticColumns := len(p.Metadata.Static)

	p.formula = fmt.Sprintf("(%d * (%d - %d - %d) + %d)", p.Metadata.Rows, allColumns, primaryColumns, staticColumns, staticColumns)

	p.result = p.Metadata.Rows*(allColumns-primaryColumns-staticColumns) + staticColumns
}

// GetResult get NOV result
func (p *NOV) GetResult() int {
	return p.result
}

// GetFormula get NOV formula
func (p *NOV) GetFormula() string {
	return p.formula
}

// PDS partition disk size data
type PDS struct {
	Metadata Metadata
	NOV      int

	formula string
	result  int
}

// Calculate calculate partition disk size
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

// GetResult get PDS result
func (p *PDS) GetResult() int {
	return p.result
}

// GetFormula get PDS formuula
func (p *PDS) GetFormula() string {
	return p.formula
}
