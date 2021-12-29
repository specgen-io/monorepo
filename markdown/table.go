package markdown

import (
	"bytes"
	"strings"
)

type Table struct {
	Header []string
	Align  []Align
	Rows   [][]string
}

func NewTable() *Table {
	return &Table{[]string{}, []Align{}, [][]string{}}
}

type Align string

const (
	Left   Align = "left"
	Center Align = "center"
	Right  Align = "right"
)

func (table *Table) AddColumn(name string, align Align) {
	table.Header = append(table.Header, name)
	table.Align = append(table.Align, align)
}

func (table *Table) AddRow(row ...string) {
	if len(row) != len(table.Header) {
		panic("Number of cell has to be equal numbers of columns")
	}
	table.Rows = append(table.Rows, row)
}

func (table *Table) ColWidth(colIndex int) int {
	width := 3
	headerWidth := len(table.Header[colIndex])
	if headerWidth > width {
		width = headerWidth
	}
	for _, row := range table.Rows {
		rowWidth := len(row[colIndex])
		if rowWidth > width {
			width = rowWidth
		}
	}
	return width
}

func space(value string, length int, spacer string) string {
	return value + strings.Repeat(spacer, length-len(value))
}

func tableRow(widths []int, row []string) string {
	cells := []string{}
	for colIndex, cell := range row {
		cells = append(cells, " "+space(cell, widths[colIndex], " ")+" ")
	}
	return "| " + strings.Join(cells, "|") + " |\n"
}

func tableAlign(widths []int, aligns []Align) string {
	cells := []string{}
	for colIndex, align := range aligns {
		cell := ""
		switch align {
		case Left:
			cell = ":" + strings.Repeat("-", widths[colIndex]-1)
		case Right:
			cell = strings.Repeat("-", widths[colIndex]-1) + ":"
		case Center:
			cell = ":" + strings.Repeat("-", widths[colIndex]-2) + ":"
		}
		cells = append(cells, " "+cell+" ")
	}
	return "| " + strings.Join(cells, "|") + " |\n"
}

func (table *Table) String() string {
	widths := make([]int, len(table.Header))
	for colIndex, _ := range table.Header {
		widths[colIndex] = table.ColWidth(colIndex)
	}
	buf := new(bytes.Buffer)
	buf.WriteString(tableRow(widths, table.Header))
	buf.WriteString(tableAlign(widths, table.Align))
	for _, row := range table.Rows {
		buf.WriteString(tableRow(widths, row))
	}
	return buf.String()
}
