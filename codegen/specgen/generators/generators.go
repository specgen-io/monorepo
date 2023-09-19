package generators

import (
	"generator"
	golang "golang/generators"
	java "java/generators"
	kotlin "kotlin/generators"
	"openapi"
	typescript "typescript/generators"
)

var All = []generator.Generator{
	golang.Models,
	golang.Client,
	golang.Service,
	java.Models,
	java.Client,
	java.Service,
	kotlin.Models,
	kotlin.Client,
	kotlin.Service,
	typescript.Models,
	typescript.Client,
	typescript.Service,
	openapi.Openapi,
}
