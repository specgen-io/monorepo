package client

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
)

func clientException(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

public class ClientException extends RuntimeException {
	public ClientException() {
		super();
	}

	public ClientException(String message) {
		super(message);
	}

	public ClientException(String message, Throwable cause) {
		super(message, cause);
	}

	public ClientException(Throwable cause) {
		super(cause);
	}
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ClientException.java"),
		Content: strings.TrimSpace(code),
	}
}
