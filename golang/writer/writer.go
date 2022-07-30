package writer

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
)

var GoConfig = generator.Config{"\t", 2, map[string]string{"PERCENT_": "%"}}

func NewGoWriter() *generator.Writer {
	return generator.NewWriter(GoConfig)
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
