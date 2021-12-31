package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

func generateClientException(thePackage Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("ClientException.java"),
		Content: strings.TrimSpace(code),
	}
}
