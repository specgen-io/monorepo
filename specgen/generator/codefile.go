package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type CodeFile struct {
	Path    string
	Content string
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func WriteFile(file *CodeFile, overwrite bool) (bool, error) {
	if overwrite || !exists(file.Path) {
		data := []byte(file.Content)

		dir := filepath.Dir(file.Path)
		_ = os.MkdirAll(dir, os.ModePerm)

		err := ioutil.WriteFile(file.Path, data, 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}
}

type ProcessedFileCallback func(wrote bool, fullpath string)

func WriteFiles(files []CodeFile, overwrite bool, processedFile ProcessedFileCallback) error {
	for _, file := range files {
		fullpath, err := filepath.Abs(file.Path)
		if err != nil {
			return err
		}
		wrote, err := WriteFile(&file, overwrite)
		if err != nil {
			return err
		}
		if processedFile != nil {
			processedFile(wrote, fullpath)
		}
	}
	return nil
}
