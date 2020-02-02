package static

import (
	"bytes"
	rice "github.com/GeertJohan/go.rice"
	"os"
	"path/filepath"
	"specgen/gen"
	"text/template"
)

func renderTemplate(content string, data interface{}) (string, error) {
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

func renderTemplateFile(box *rice.Box, template string, outputPath string, data interface{}) (*gen.TextFile, error) {
	templateContent, err := box.String(template)
	if err != nil {
		return nil, err
	}

	text, err := renderTemplate(templateContent, data)
	if err != nil {
		return nil, err
	}

	file := gen.TextFile{Path: outputPath, Content: text}
	return &file, nil
}

func RenderTemplate(template string, outPath string, data interface{}) ([]gen.TextFile, error) {
	riceBox := rice.MustFindBox("code")

	files := []gen.TextFile{}

	err := riceBox.Walk(template, func(itemTemplatePath string, fileInfo os.FileInfo, walkError error) error {
		itemOutputPath, walkError := renderTemplate(itemTemplatePath[len(template):], data)
		if walkError != nil {
			return walkError
		}

		outputPath := filepath.Join(outPath, itemOutputPath)
		if fileInfo.IsDir() {
			return os.MkdirAll(outputPath, os.ModePerm)
		} else {
			file, internalError := renderTemplateFile(riceBox, itemTemplatePath, outputPath, data)
			if internalError != nil {
				return internalError
			}
			files = append(files, *file)
			return nil
		}
	})
	return files, err
}