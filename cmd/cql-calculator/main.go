package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	calculator "github.com/johnnywidth/cql-calculator"
	"gopkg.in/yaml.v2"
)

// MEGABYTE megabyte
const MEGABYTE = 1.0 << (10 * 2)

func main() {
	fileName := flag.String("file", "", "yml file for save cql parser result")
	query := flag.String("query", "", "CQL query")
	flag.Parse()

	if *query == "" && *fileName == "" {
		flag.PrintDefaults()
		return
	}

	meta := &calculator.Metadata{}

	if *query != "" {
		fmt.Print("Enter rows count per one partition: ")

		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			panic(err)
		}

		meta, err = calculator.PopulateTableMetadata(*query, i)
		if err != nil {
			panic(err)
		}

		if len(meta.GetNotSpecifiedSizes()) > 0 {
			for _, v := range meta.GetNotSpecifiedSizes() {
				fmt.Printf("Enter (avarage) size for `%s (%s)` column: ", v.Name, v.Type)

				var i int
				_, err = fmt.Scanf("%d", &i)
				if err != nil {
					panic(err)
				}
				v.Size = i

				meta.SpecifyCustomSize(v)
			}
		}

		writeMetaToFile(meta, *fileName)
	}

	populateMetaFromFile(meta, *fileName, *query)

	nov := calculator.NOV{
		Metadata: meta,
	}
	nov.Calculate()

	pds := calculator.PDS{
		Metadata: meta,
		NOV:      nov.GetResult(),
	}
	pds.Calculate()

	fmt.Printf("Number of Values:\n%s = %d\n\n", nov.GetFormula(), nov.GetResult())
	fmt.Printf("Partition Size on Disk:\n%s = %d bytes (%.2f Mb)\n", pds.GetFormula(), pds.GetResult(), float64(pds.GetResult())/MEGABYTE)
}

func writeMetaToFile(meta *calculator.Metadata, fileName string) {
	if fileName != "" {
		data, err := yaml.Marshal(meta)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fileName, data, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func populateMetaFromFile(meta *calculator.Metadata, fileName, query string) {
	if fileName != "" && query == "" {
		f, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(f, &meta)
		if err != nil {
			panic(err)
		}
	}
}
