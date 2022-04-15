package client

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/spec"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func responseBodyName(response *spec.NamedResponse) string {
	return fmt.Sprintf(`responseBody%s`, casee.ToPascalCase(spec.HttpStatusName(spec.HttpStatusCode(response.Name))))
}
