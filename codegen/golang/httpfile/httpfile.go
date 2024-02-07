package httpfile

import (
	"generator"
	"golang/module"
	"golang/writer"
)

func GenerateFile(httpFileModule module.Module) *generator.CodeFile {
	w := writer.New(httpFileModule, `file.go`)
	w.Lines(`
import (
	"io"
	"os"
	"strings"
)

type File struct {
	Name    string
	Content io.ReadCloser
}

func NewFile(fileName, fileContent string) *File {
	content := io.NopCloser(strings.NewReader(fileContent))

	return &File{Name: fileName, Content: content}
}

func NewFileFromFS(fileName string) (*File, error) {
	fileContent, err := os.Open(fileName)

	return &File{Name: fileName, Content: fileContent}, err
}
`)
	return w.ToCodeFile()
}
