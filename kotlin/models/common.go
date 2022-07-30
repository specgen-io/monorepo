package models

import (
	"github.com/specgen-io/specgen/spec/v2"
)

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}
