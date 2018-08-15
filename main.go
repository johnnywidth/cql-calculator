package main

import (
	"flag"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

const MEGABYTE = 1.0 << (10 * 2)

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

var m Metadata

func main() {
	fileName := flag.String("file", "", "")
	generate := flag.String("generate", "", "")
	flag.Parse()

	if *generate != "" {
		generateFromCQL(*generate)
		return
	}

	f, err := ioutil.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}

	m = Metadata{}
	err = yaml.Unmarshal(f, &m)
	if err != nil {
		panic(err)
	}

	nov := numberOfValues(m)
	log.Printf("Number of Values:\n%d\n\n", nov)

	pds := partitionDiskSize(m, nov)
	log.Printf("Partition Size on Disk:\n%d bytes\n%.2f Mb", pds, float64(pds)/MEGABYTE)
}

func numberOfValues(t Metadata) int {
	allColumns := len(t.Column) + len(t.Partition) + len(t.Cluster) + len(t.Static)
	primaryColumns := len(t.Partition) + len(t.Cluster)
	staticColumns := len(t.Static)

	log.Printf("Number of Values(%d*(%d-%d-%d) + %d)", t.Rows, allColumns, primaryColumns, staticColumns, staticColumns)

	return t.Rows*(allColumns-primaryColumns-staticColumns) + staticColumns
}

func partitionDiskSize(t Metadata, nov int) int {
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

	log.Printf("Partition Size on Disk(%d + %d + (%d * %d) + (8 * %d))", sofPK, sofS, t.Rows, ek, nov)

	return sofPK + sofS + (t.Rows * ek) + (8 * nov)
}
