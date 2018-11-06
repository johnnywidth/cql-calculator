package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/johnnywidth/cql-calculator"
	"gopkg.in/yaml.v2"
)

// MEGABYTE megabyte
const MEGABYTE = 1.0 << (10 * 2)

func main() {
	fileName := flag.String("file", "", "")
	generate := flag.String("generate", "", "")
	flag.Parse()

	if *fileName != "" && *generate != "" {
		fmt.Print("Enter rows count per one partition: ")

		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			panic(err)
		}

		meta, err := calculator.PopulateTableMetadata(*generate, i)
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

				err = meta.SpecifyCustomSize(v)
				if err != nil {
					panic(err)
				}
			}
		}

		data, err := yaml.Marshal(meta)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(*fileName, data, 0755)
		if err != nil {
			panic(err)
		}
	}

	f, err := ioutil.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}

	m := calculator.Metadata{}
	err = yaml.Unmarshal(f, &m)
	if err != nil {
		panic(err)
	}

	nov := calculator.NOV{
		Metadata: m,
	}
	nov.Calculate()

	pds := calculator.PDS{
		Metadata: m,
		NOV:      nov.GetResult(),
	}
	pds.Calculate()

	fmt.Printf("Number of Values:\n%s = %d\n\n", nov.GetFormula(), nov.GetResult())
	fmt.Printf("Partition Size on Disk:\n%s = %d bytes (%.2f Mb)\n", pds.GetFormula(), pds.GetResult(), float64(pds.GetResult())/MEGABYTE)
}
