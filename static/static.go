package static

import rice "github.com/GeertJohan/go.rice"

func StaticCode(path string) (string, error) {
	box := rice.MustFindBox("code")
	return box.String(path)
}