package client

import (
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

func clientException(thePackage modules.Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("ClientException.kt"),
		Content: strings.TrimSpace(code),
	}
}
