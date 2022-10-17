package writer

import (
	"generator"
	"golang/module"
	"strings"
)

func GoConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{"PERCENT_": "%"}}
}

func New(module module.Module, filename string) generator.Writer {
	config := GoConfig()
	w := generator.NewWriter2(module.GetPath(filename), config)
	w.Line("package %s", module.Name)
	w.EmptyLine()
	return w
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

func WriteAlignedLines(w generator.Writer, lines [][]string) {
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
