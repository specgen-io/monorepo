package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateClientException(thePackage Module) *gen.TextFile {
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

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("ClientException.java"),
		Content: strings.TrimSpace(code),
	}
}
