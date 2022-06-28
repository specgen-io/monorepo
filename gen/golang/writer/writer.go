package writer

import (
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

func NewGoWriter() *generator.Writer {
	return generator.NewWriter("\t", 2)
}

func colWidth(lines [][]string, colIndex int) int {
	width := 0
	for _, line := range lines {
		rowWidth := len(line[colIndex])
		if rowWidth > width {
			width = rowWidth
		}
	}
	return width
}

func space(value string, length int, spacer string) string {
	return value + strings.Repeat(spacer, length-len(value))
}

func WriteAlignedLines(w *generator.Writer, lines [][]string) {
	widths := make([]int, len(lines[0]))
	for colIndex, _ := range lines[0] {
		widths[colIndex] = colWidth(lines, colIndex)
	}
	for _, line := range lines {
		lineStr := ""
		for colIndex, cell := range line {
			lineStr += cell
			if colIndex != len(line)-1 {
				lineStr += strings.Repeat(" ", widths[colIndex]-len(cell)) + " "
			}
		}
		w.Line(lineStr)
	}
}
