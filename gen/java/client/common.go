package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}
