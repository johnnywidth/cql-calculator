package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/johnnywidth/cql-calculator/src"
	yaml "gopkg.in/yaml.v2"
)

const MEGABYTE = 1.0 << (10 * 2)

var m src.Metadata

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

	m = src.Metadata{}
	err = yaml.Unmarshal(f, &m)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	nov := src.NumberOfValues(m)
	fmt.Printf("Number of Values:\n%d\n\n", nov)

	pds := src.PartitionDiskSize(m, nov)
	fmt.Printf("Partition Size on Disk:\n%d bytes\n%.2f Mb\n", pds, float64(pds)/MEGABYTE)
}
