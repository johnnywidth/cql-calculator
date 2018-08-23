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

func NumberOfValues(t Metadata) int {
	allColumns := len(t.Column) + len(t.Partition) + len(t.Cluster) + len(t.Static)
	primaryColumns := len(t.Partition) + len(t.Cluster)
	staticColumns := len(t.Static)

	fmt.Printf("Number of Values(%d*(%d-%d-%d) + %d)\n", t.Rows, allColumns, primaryColumns, staticColumns, staticColumns)

	return t.Rows*(allColumns-primaryColumns-staticColumns) + staticColumns
}

func PartitionDiskSize(t Metadata, nov int) int {
	var sofPK int
	for _, v := range t.Partition {
		sofPK += v.Size
	}

	var sofS int
	for _, v := range t.Static {
		sofS += v.Size
	}

	var sofCK int
	for _, v := range t.Cluster {
		sofCK += v.Size
	}

	var ek int
	for _, v := range t.Column {
		ek += v.Size + sofCK
	}

	fmt.Printf("Partition Size on Disk(%d + %d + (%d * %d) + (8 * %d))\n", sofPK, sofS, t.Rows, ek, nov)

	return sofPK + sofS + (t.Rows * ek) + (8 * nov)
}
