package main

import (
	"flag"
	"fmt"
	"io/ioutil"

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

	if *fileName != "" && *generate != "" {
		generateFromCQL(*fileName, *generate)
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

	fmt.Println()
	nov := numberOfValues(m)
	fmt.Printf("Number of Values:\n%d\n\n", nov)

	pds := partitionDiskSize(m, nov)
	fmt.Printf("Partition Size on Disk:\n%d bytes\n%.2f Mb\n", pds, float64(pds)/MEGABYTE)
}
