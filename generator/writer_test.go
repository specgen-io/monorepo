package generator

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Basic_Lines(t *testing.T) {
	expected := `
line1
line2
`
	w := NewWriter(Config{"", 0, nil})
	w.Line("line1")
	w.Line("line2")
	assert.Equal(t, strings.TrimSpace(w.String()), strings.TrimSpace(expected))
}

func Test_Indentation(t *testing.T) {
	expected := `
line1
  line2
    line3
`
	w := NewWriter(Config{"  ", 0, nil})
	w.Line("line1")
	w.Indent()
	w.Line("line2")
	w.Indent()
	w.Line("line3")
	assert.Equal(t, strings.TrimSpace(w.String()), strings.TrimSpace(expected))
}

func Test_Visual_Indentation(t *testing.T) {
	expected := `
line1
  line2
    line3
`
	w := NewWriter(Config{"  ", 4, nil})
	w.Line("line1")
	w.Line("    line2")
	w.Line("        line3")
	assert.Equal(t, strings.TrimSpace(w.String()), strings.TrimSpace(expected))
}
