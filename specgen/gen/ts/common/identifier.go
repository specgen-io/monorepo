package common

import (
	"fmt"
	"regexp"
)

func checkFormat(regex string, value string) bool {
	isMatching, err := regexp.MatchString(regex, value)
	if err != nil {
		panic(err)
	}
	return isMatching
}

var tsIdentifierFormat = "^[a-zA-Z_]([a-zA-Z0-9_])*$"

func TSIdentifier(name string) string {
	if !checkFormat(tsIdentifierFormat, name) {
		return fmt.Sprintf("'%s'", name)
	}
	return name
}
