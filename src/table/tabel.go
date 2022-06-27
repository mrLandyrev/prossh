package table

import (
	"fmt"
)

type Table struct {
	Headers []string
	Data    [][]string
}

func (table *Table) AddColumn(name string) (err error) {
	if len(table.Data) > 0 {
		err = fmt.Errorf("Add column is unavailable: table has data")
		return
	}
	table.Headers = append(table.Headers, name)
	return
}

func (table *Table) AddRow(row []string) (err error) {
	if len(table.Headers) != len(row) {
		err = fmt.Errorf("Add row is unavailable: len is incorrect")
		return
	}
	table.Data = append(table.Data, row)
	return
}

func (table Table) String() (res string) {
	res = ""

	resSize := make([]int, len(table.Headers))
	for colIndex := range table.Headers {
		resSize[colIndex] = len(table.Headers[colIndex])
		for rowIndex := range table.Data {
			if len(table.Data[rowIndex][colIndex]) > resSize[colIndex] {
				resSize[colIndex] = len(table.Data[rowIndex][colIndex])
			}
		}
	}

	for colIndex := range table.Headers {
		res += appendSpacesToWidth(table.Headers[colIndex], resSize[colIndex])
		res += " "
	}

	if len(table.Data) == 0 {
		res += "\nEmpty"
	} else {
		for rowIndex := range table.Data {
			res += "\n"
			for colIndex := range table.Data[rowIndex] {
				res += appendSpacesToWidth(table.Data[rowIndex][colIndex], resSize[colIndex])
				res += " "
			}
		}
	}
	res += "\n"

	return
}

func appendSpacesToWidth(str string, width int) (res string) {
	res = str
	for len(res) < width {
		res += " "
	}
	return
}
