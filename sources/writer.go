package sources

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Writer struct {
	buffer                    *bytes.Buffer
	indentationStr            string
	leadSpacesIndentationSize int

	indentation int
}

func NewWriter(indentationStr string, leadSpacesIndentationSize int) *Writer {
	return &Writer{new(bytes.Buffer), indentationStr, leadSpacesIndentationSize, 0}
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
	if writer.leadSpacesIndentationSize > 0 {
		prefix := strings.Repeat(" ", writer.leadSpacesIndentationSize)
		line, indentation = trimPrefix(line, prefix)
	}
	realIndentation := indentation + writer.indentation
	indentationStr := strings.Repeat(writer.indentationStr, realIndentation)
	line = indentationStr + line + "\n"
	io.WriteString(writer.buffer, line)
}

func (writer *Writer) EmptyLine() {
	io.WriteString(writer.buffer, "\n")
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
		writer.buffer,
		writer.indentationStr,
		writer.leadSpacesIndentationSize,
		writer.indentation + 1,
	}
}

func (writer *Writer) IndentedWith(size int) *Writer {
	return &Writer{
		writer.buffer,
		writer.indentationStr,
		writer.leadSpacesIndentationSize,
		writer.indentation + size,
	}
}

func (writer *Writer) String() string {
	return writer.buffer.String()
}
