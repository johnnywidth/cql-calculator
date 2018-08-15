package main

func GetSizeByType(t string) int {
	switch t {
	case "tinyint":
		return 1
	case "int":
		return 4
	case "timestamp":
		return 8
	case "uuid":
		return 16
	}

	return -1
}
