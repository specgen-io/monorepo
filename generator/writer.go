package generator

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Config struct {
	IndentationStr            string
	LeadSpacesIndentationSize int
	Substitutions             map[string]string
}

type Writer struct {
	filename    string
	config      Config
	buffer      *bytes.Buffer
	indentation int
}

func NewWriter(config Config) *Writer {
	return &Writer{
		"",
		config,
		new(bytes.Buffer),
		0,
	}
}

func NewWriter2(filename string, config Config) *Writer {
	return &Writer{
		filename,
		config,
		new(bytes.Buffer),
		0,
	}
}

func (writer *Writer) ToCodeFile() *CodeFile {
	return &CodeFile{
		Path:    writer.filename,
		Content: writer.String(),
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

func (writer Writer) Write(s string) {
	io.WriteString(writer.buffer, substitute(s, writer.config.Substitutions))
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

func (writer *Writer) Line(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	indentation := 0
	if writer.config.LeadSpacesIndentationSize > 0 {
		prefix := strings.Repeat(" ", writer.config.LeadSpacesIndentationSize)
		line, indentation = trimPrefix(line, prefix)
	}
	realIndentation := indentation + writer.indentation
	indentationStr := strings.Repeat(writer.config.IndentationStr, realIndentation)
	line = indentationStr + line + "\n"
	writer.Write(line)
}

func (writer *Writer) Lines(format string, args ...interface{}) {
	code := fmt.Sprintf(strings.Trim(format, "\n"), args...)
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		writer.Line(line)
	}
}

func (writer *Writer) EmptyLine() {
	writer.Write("\n")
}

func (writer *Writer) Indent() {
	writer.indentation = writer.indentation + 1
}

func (writer *Writer) Unindent() {
	writer.indentation = writer.indentation - 1
}

func (writer *Writer) IndentWith(size int) {
	writer.indentation = writer.indentation + size
}

func (writer *Writer) UnindentWith(size int) {
	writer.indentation = writer.indentation - size
}

func (writer *Writer) Indented() *Writer {
	return &Writer{
		writer.filename,
		writer.config,
		writer.buffer,
		writer.indentation + 1,
	}
}

func (writer *Writer) IndentedWith(size int) *Writer {
	return &Writer{
		writer.filename,
		writer.config,
		writer.buffer,
		writer.indentation + size,
	}
}

func (writer *Writer) String() string {
	return writer.buffer.String()
}
