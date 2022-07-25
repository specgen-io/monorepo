package client

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
)

func clientException(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import java.lang.RuntimeException

class ClientException : RuntimeException {
    constructor() : super()
    constructor(message: String) : super(message)
    constructor(message: String, cause: Throwable) : super(message, cause)
    constructor(cause: Throwable) : super(cause)
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ClientException.kt"),
		Content: strings.TrimSpace(code),
	}
}
