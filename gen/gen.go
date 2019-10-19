package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type TextFile struct {
	Path    string
	Content string
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func WriteFile(file *TextFile, overwrite bool) error {
	if overwrite || !Exists(file.Path) {
		data := []byte(file.Content)

		dir := filepath.Dir(file.Path)
		_ = os.MkdirAll(dir, os.ModePerm)

		return ioutil.WriteFile(file.Path, data, 0644)
	}
	return nil
}

func WriteFiles(files []TextFile, overwrite bool) error {
	for _, file := range files {
		err := WriteFile(&file, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}
