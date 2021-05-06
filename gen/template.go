package gen

import (
	"bytes"
	"text/template"
)

func ExecuteTemplate(content string, data interface{}) (string, error) {
	var buffer bytes.Buffer

	t, err := template.New("template").Delims("[[", "]]").Parse(content)
	if err != nil {
		return "", err
	}

	err = t.ExecuteTemplate(&buffer, "template", data)
	if err != nil {
		return "", err
	}
	text := buffer.String()

	return text, nil
}