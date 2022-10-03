package client

import (
	"strings"

	"generator"
)

func (g *Generator) Exceptions() *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{g.Packages.Root.PackageName})
	return &generator.CodeFile{
		Path:    g.Packages.Root.GetPath("ClientException.java"),
		Content: strings.TrimSpace(code),
	}
}
