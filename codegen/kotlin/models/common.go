package models

import (
	"spec"
)

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}
