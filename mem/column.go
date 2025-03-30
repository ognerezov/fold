package mem

import (
	"fmt"
	"fold/openapi"
	"google.golang.org/api/sheets/v4"
)

type ColumnDefinition struct {
	name          string
	isIndex       bool
	isUnique      bool
	foreignTable  string
	foreignColumn string
	foreignUnique bool
	dataType      string
	number        int
}

func (c *ColumnDefinition) ToString() string {
	if c.isIndex {
		return fmt.Sprintf("[%s*]", c.name)
	}
	if c.isIndex {
		return fmt.Sprintf("%s*", c.name)
	}
	return c.name
}

func (c *ColumnDefinition) Name() string {
	return c.name
}

func (c *ColumnDefinition) Number() int {
	return c.number
}

func (c *ColumnDefinition) SetDataType(t string) {
	if c.dataType == t {
		return
	}
	if c.dataType == "" {
		c.dataType = t
		return
	}
	c.dataType = String
}

func (c *ColumnDefinition) Schema() openapi.Schema {
	switch c.dataType {
	case Int:
	case Float:
		return openapi.Schema{
			Type: "number",
		}
	case Array:
		return openapi.Schema{
			Type: "array",
		}
	case Struct:
		return openapi.Schema{
			Type: "object",
		}
	case Bool:
		return openapi.Schema{
			Type: "boolean",
		}
	default:
		return openapi.Schema{
			Type: String,
		}
	}
	return openapi.Schema{
		Type: String,
	}
}

func (c *ColumnDefinition) CellData() *sheets.CellData {
	v := sheets.ExtendedValue{
		StringValue: &c.name,
	}
	return &sheets.CellData{
		UserEnteredValue: &v,
	}
}

func ColumnsPrintln(columns []*ColumnDefinition) {
	fmt.Print("_ |")
	for _, column := range columns {
		fmt.Print(column.ToString() + " |")
	}
	fmt.Println()
}

func SimpleDefinition(name string, isIndex bool, index int) *ColumnDefinition {
	return &ColumnDefinition{name: name, isIndex: isIndex, isUnique: isIndex, number: index}
}

func (c *ColumnDefinition) IsIndex() bool {
	return c.isIndex
}

func (c *ColumnDefinition) BackReference(tableName string, originalColNumber int) *ColumnDefinition {
	return &ColumnDefinition{
		name:          fmt.Sprintf("%s_%s", tableName, c.name),
		isIndex:       c.isIndex,
		isUnique:      c.isUnique,
		dataType:      c.dataType,
		number:        originalColNumber,
		foreignTable:  tableName,
		foreignColumn: c.name,
		foreignUnique: c.isUnique,
	}
}
