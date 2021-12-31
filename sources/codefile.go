package sources

import (
	"github.com/specgen-io/specgen/v2/console"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CodeFile struct {
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

func WriteFile(file *CodeFile, overwrite bool) error {
	fullpath, err := filepath.Abs(file.Path)
	if err != nil {
		return err
	}
	if overwrite || !Exists(file.Path) {
		data := []byte(file.Content)

		dir := filepath.Dir(file.Path)
		_ = os.MkdirAll(dir, os.ModePerm)

		console.PrintLn("Writing:", fullpath)
		return ioutil.WriteFile(file.Path, data, 0644)
	} else {
		console.PrintLn("Skipping:", fullpath)
	}
	return nil
}

func WriteFiles(files []CodeFile, overwrite bool) error {
	for _, file := range files {
		err := WriteFile(&file, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}
