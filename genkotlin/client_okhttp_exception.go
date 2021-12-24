package genkotlin

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateClientException(thePackage Module) *gen.TextFile {
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

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("ClientException.kt"),
		Content: strings.TrimSpace(code),
	}
}
