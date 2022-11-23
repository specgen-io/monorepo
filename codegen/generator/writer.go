package generator

import (
	"fmt"
	"strings"
)

type Writer interface {
	Line(format string, args ...interface{})
	Lines(format string)
	Template(data map[string]string, content string)
	EmptyLine()
	LineAligned(format string, args ...interface{})
	Indent()
	Unindent()
	IndentWith(size int)
	UnindentWith(size int)
	Indented() Writer
	IndentedWith(size int) Writer
	String() string
	Code() []string
}

type Config struct {
	IndentationStr            string
	LeadSpacesIndentationSize int
	Substitutions             map[string]string
}

type content struct {
	lines        []string
	linesAligned []string
}

func (c *content) Add(s string) {
	c.lines = append(c.lines, s)
}

type writer struct {
	config      Config
	content     *content
	indentation int
}

func NewWriter(config Config) Writer {
	return &writer{
		config,
		&content{[]string{}, []string{}},
		0,
	}
}

func substitute(s string, substitutions map[string]string) string {
	result := s
	if substitutions != nil {
		for oldStr, newStr := range substitutions {
			result = strings.Replace(result, oldStr, newStr, -1)
		}
	}
	return result
}

func wrapKeys(vars map[string]string, prefix, postfix string) map[string]string {
	result := map[string]string{}
	for key, value := range vars {
		result[fmt.Sprintf(`%s%s%s`, prefix, key, postfix)] = value
	}
	return result
}

func (w *writer) write(s string) {
	w.checkAligned()
	w.content.Add(s)
}

func (w *writer) LineAligned(format string, args ...interface{}) {
	theline := fmt.Sprintf(format, args...)
	w.content.linesAligned = append(w.content.linesAligned, theline)
}

func (w *writer) checkAligned() {
	if len(w.content.linesAligned) == 0 {
		return
	}
	lines := w.content.linesAligned
	w.content.linesAligned = []string{}

	linesParts := [][]string{}
	for _, line := range lines {
		linesParts = append(linesParts, strings.Split(line, " "))
	}

	widths := make([]int, len(linesParts[0]))
	for colIndex, _ := range linesParts[0] {
		widths[colIndex] = colWidth(linesParts, colIndex)
	}
	for _, line := range linesParts {
		lineStr := ""
		for colIndex, cell := range line {
			lineStr += cell
			if colIndex != len(line)-1 {
				lineStr += strings.Repeat(" ", widths[colIndex]-len(cell)) + " "
			}
		}
		w.line(lineStr)
	}
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

func trimPrefix(str string, prefix string) (string, int) {
	count := 0
	trimmed := str
	for strings.HasPrefix(trimmed, prefix) {
		trimmed = strings.TrimPrefix(trimmed, prefix)
		count = count + 1
	}
	return trimmed, count
}

func (w *writer) line(theline string) {
	indentation := 0
	if w.config.LeadSpacesIndentationSize > 0 {
		prefix := strings.Repeat(" ", w.config.LeadSpacesIndentationSize)
		theline, indentation = trimPrefix(theline, prefix)
	}
	realIndentation := indentation + w.indentation
	indentationStr := strings.Repeat(w.config.IndentationStr, realIndentation)
	theline = indentationStr + theline + "\n"
	theline = substitute(theline, w.config.Substitutions)
	w.write(theline)
}

func (w *writer) Line(format string, args ...interface{}) {
	theline := fmt.Sprintf(format, args...)
	w.line(theline)
}

func (w *writer) Lines(content string) {
	code := strings.Trim(content, "\n")
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		w.line(line)
	}
}

func (w *writer) Template(data map[string]string, content string) {
	code := substitute(content, wrapKeys(data, "[[.", "]]"))
	w.Lines(code)
}

func (w *writer) EmptyLine() {
	w.write("\n")
}

func (w *writer) setIndentation(value int) {
	w.checkAligned()
	w.indentation = value
}

func (w *writer) Indent() {
	w.setIndentation(w.indentation + 1)
}

func (w *writer) Unindent() {
	w.setIndentation(w.indentation - 1)
}

func (w *writer) IndentWith(size int) {
	w.setIndentation(w.indentation + size)
}

func (w *writer) UnindentWith(size int) {
	w.setIndentation(w.indentation - size)
}

func (w *writer) Indented() Writer {
	return w.IndentedWith(1)
}

func (w *writer) IndentedWith(size int) Writer {
	w.checkAligned()
	return &writer{
		w.config,
		w.content,
		w.indentation + size,
	}
}

func (w *writer) String() string {
	w.checkAligned()
	return strings.Join(w.content.lines, ``)
}

func (w *writer) Code() []string {
	w.checkAligned()
	return w.content.lines
}
