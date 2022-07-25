package client

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"strings"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}

func trimSlash(param string) string {
	return strings.Trim(param, "/")
}
