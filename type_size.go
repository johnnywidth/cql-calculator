package calculator

import (
	"errors"

	cql "github.com/johnnywidth/cql-calculator/cql-parser"
)

// CustomSize size which need custom specify
type CustomSize struct {
	Name string
	Type string
	Size int
}

// GetNotSpecifiedSizes get all not specified sizes
func (m *Metadata) GetNotSpecifiedSizes() map[string]CustomSize {
	return m.customSizes
}

// SpecifyCustomSize specify custom size
func (m *Metadata) SpecifyCustomSize(s CustomSize) error {
	_, ok := m.customSizes[s.Name]
	if !ok || s.Size <= 0 {
		return errors.New("invalid size")
	}

	for k, v := range m.Partition {
		if v.Name == s.Name {
			m.Partition[k].Size = s.Size
			return nil
		}
	}

	for k, v := range m.Cluster {
		if v.Name == s.Name {
			m.Cluster[k].Size = s.Size
			return nil
		}
	}

	for k, v := range m.Static {
		if v.Name == s.Name {
			m.Static[k].Size = s.Size
			return nil
		}
	}

	for k, v := range m.Column {
		if v.Name == s.Name {
			m.Column[k].Size = s.Size
			return nil
		}
	}

	return errors.New("did not specify")
}

func (m *Metadata) addCustomSize(c cql.Column) {
	if m.customSizes == nil {
		m.customSizes = make(map[string]CustomSize, 0)
	}
	m.customSizes[c.Name] = CustomSize{
		Name: c.Name,
		Type: c.Type,
	}
}

// GetSizeByType get zise by type
func GetSizeByType(n, t string) (int, error) {
	switch t {
	case "decimal", "duration":
		return -1, nil
	case "boolean", "tinyint", "smallint":
		return 1, nil
	case "date", "int", "float", "inet":
		return 4, nil
	case "bigint", "counter", "time", "timestamp", "double", "varint":
		return 8, nil
	case "uuid", "timeuuid":
		return 16, nil
	case "ascii", "text", "varchar", "blob", "map", "list", "set":
		return -1, errors.New("specify size")
	}

	return -1, errors.New("type not exist in cql")
}
