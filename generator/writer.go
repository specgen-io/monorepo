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
	ToCodeFile() *CodeFile
	String() string
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

type TheWriter struct {
	filename    string
	config      Config
	content     *content
	indentation int
}

func NewWriter(filename string, config Config) Writer {
	return &TheWriter{
		filename,
		config,
		&content{[]string{}, []string{}},
		0,
	}
}

func (w *TheWriter) ToCodeFile() *CodeFile {
	return &CodeFile{
		Path:    w.filename,
		Content: w.String(),
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

func (w *TheWriter) write(s string) {
	w.checkAligned()
	w.content.Add(s)
}

func (w *TheWriter) LineAligned(format string, args ...interface{}) {
	theline := fmt.Sprintf(format, args...)
	w.content.linesAligned = append(w.content.linesAligned, theline)
}

func (w *TheWriter) checkAligned() {
	if len(w.content.linesAligned) == 0 {
		return
	}
	lines := w.content.linesAligned
	w.content.linesAligned = []string{}

	linesParts := [][]string{}
	for _, line := range lines {
		linesParts = append(linesParts, strings.Split(line, " "))
	}

	widths := make([]int, len(lines[0]))
	for colIndex, _ := range lines[0] {
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

func (w *TheWriter) line(theline string) {
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

func (w *TheWriter) Line(format string, args ...interface{}) {
	theline := fmt.Sprintf(format, args...)
	w.line(theline)
}

func (w *TheWriter) Lines(content string) {
	code := strings.Trim(content, "\n")
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		w.line(line)
	}
}

func (w *TheWriter) Template(data map[string]string, content string) {
	code := substitute(content, wrapKeys(data, "[[.", "]]"))
	w.Lines(code)
}

func (w *TheWriter) EmptyLine() {
	w.write("\n")
}

func (w *TheWriter) Indent() {
	w.indentation = w.indentation + 1
}

func (w *TheWriter) Unindent() {
	w.indentation = w.indentation - 1
}

func (w *TheWriter) IndentWith(size int) {
	w.indentation = w.indentation + size
}

func (w *TheWriter) UnindentWith(size int) {
	w.indentation = w.indentation - size
}

func (w *TheWriter) Indented() Writer {
	return &TheWriter{
		w.filename,
		w.config,
		w.content,
		w.indentation + 1,
	}
}

func (w *TheWriter) IndentedWith(size int) Writer {
	return &TheWriter{
		w.filename,
		w.config,
		w.content,
		w.indentation + size,
	}
}

func (w *TheWriter) String() string {
	w.checkAligned()
	return strings.Join(w.content.lines, ``)
}
