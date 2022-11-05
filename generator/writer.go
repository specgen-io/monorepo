package generator

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Writer interface {
	Line(format string, args ...interface{})
	Lines(format string)
	Template(data map[string]string, content string)
	EmptyLine()
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

type writer struct {
	filename    string
	config      Config
	buffer      *bytes.Buffer
	indentation int
}

func NewWriter(filename string, config Config) Writer {
	return &writer{
		filename,
		config,
		new(bytes.Buffer),
		0,
	}
}

func (w *writer) ToCodeFile() *CodeFile {
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

func (w *writer) write(s string) {
	io.WriteString(w.buffer, substitute(s, w.config.Substitutions))
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

func (w *writer) Indent() {
	w.indentation = w.indentation + 1
}

func (w *writer) Unindent() {
	w.indentation = w.indentation - 1
}

func (w *writer) IndentWith(size int) {
	w.indentation = w.indentation + size
}

func (w *writer) UnindentWith(size int) {
	w.indentation = w.indentation - size
}

func (w *writer) Indented() Writer {
	return &writer{
		w.filename,
		w.config,
		w.buffer,
		w.indentation + 1,
	}
}

func (w *writer) IndentedWith(size int) Writer {
	return &writer{
		w.filename,
		w.config,
		w.buffer,
		w.indentation + size,
	}
}

func (w *writer) String() string {
	return w.buffer.String()
}
