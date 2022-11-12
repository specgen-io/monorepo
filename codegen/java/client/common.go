package client

import (
	"fmt"
	"spec"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}
