package client

import (
	"generator"
	"java/writer"
)

func (g *Generator) Exceptions() *generator.CodeFile {
	w := writer.New(g.Packages.Root, `ClientException`)
	w.Lines(`
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
`)
	return w.ToCodeFile()
}
