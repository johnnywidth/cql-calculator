package main

import "fmt"

func GetSizeByType(n, t string) int {
	switch t {
	case "tinyint":
		return 1
	case "int":
		return 4
	case "timestamp":
		return 8
	case "uuid":
		return 16
	case "string":
		fmt.Printf("Enter size (avarage) for `%s (%s)` type: ", n, t)
		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			panic(err)
		}
		return i
	}

	return -1
}
