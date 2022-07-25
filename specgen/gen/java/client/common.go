package client

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}
