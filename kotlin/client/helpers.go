package client

import (
	"fmt"
	"spec"
	"strings"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func trimSlash(param string) string {
	return strings.Trim(param, "/")
}
