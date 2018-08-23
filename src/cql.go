package src

import "fmt"

func GetSizeByType(n, t string) int {
	switch t {
	case "decimal", "duration":
		return -1
	case "boolean", "tinyint", "smallint":
		return 1
	case "date", "int", "float", "inet":
		return 4
	case "bigint", "counter", "time", "timestamp", "double", "varint":
		return 8
	case "uuid", "timeuuid":
		return 16
	case "ascii", "text", "varchar", "blob", "map", "list", "set":
		fmt.Printf("Enter (avarage) size for `%s (%s)` column: ", n, t)
		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			panic(err)
		}
		return i
	}

	return -1
}
