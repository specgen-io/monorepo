package markdown

import (
	"bytes"
	"fmt"
)

type Markdown struct {
	Buffer *bytes.Buffer
}

func NewMarkdown() *Markdown {
	return &Markdown{new(bytes.Buffer)}
}

func (markdown *Markdown) String() string {
	return markdown.Buffer.String()
}

func (markdown *Markdown) Header1(value string) {
	markdown.Buffer.WriteString("\n# " + value + "\n")
}

func (markdown *Markdown) Header2(value string) {
	markdown.Buffer.WriteString("\n## " + value + "\n")
}

func (markdown *Markdown) Paragraph(value string) {
	markdown.Buffer.WriteString("\n" + value + "\n")
}

func (markdown *Markdown) CodeBlock(value string) {
	markdown.Buffer.WriteString(fmt.Sprintf("\n```\n%s\n```\n", value))
}

func (markdown *Markdown) ListItem(value string) {
	markdown.Buffer.WriteString("* " + value + "\n")
}

func (markdown *Markdown) Line(value string) {
	markdown.Buffer.WriteString(value + "\n")
}

func (markdown *Markdown) Bold(value string) {
	markdown.Buffer.WriteString("**" + value + "**")
}

func (markdown *Markdown) EmptyLine() {
	markdown.Buffer.WriteString("\n")
}
