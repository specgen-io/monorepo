package scala

import (
	"strings"
)

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}
