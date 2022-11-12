package client

import (
	"fmt"
	"generator"
	"java/writer"
	"spec"
)

func (g *Generator) Exceptions(errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.clientException())
	for _, errorResponse := range *errors {
		files = append(files, *g.inheritedClientException(&errorResponse))
	}

	return files
}

func (g *Generator) clientException() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ClientException`)
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

func (g *Generator) inheritedClientException(error *spec.Response) *generator.CodeFile {
	errorName := g.Types.Java(&error.Type.Definition)
	className := fmt.Sprintf(`%sException`, errorName)
	w := writer.New(g.Packages.Errors, className)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
			`ErrorName`:           errorName,
		}, `
import [[.ErrorsModelsPackage]].*;

public class [[.ClassName]] extends ClientException {
	private final [[.ErrorName]] error;

	public [[.ClassName]]([[.ErrorName]] error) {
		super("Body: %s" + error);
		this.error = error;
	}

	public [[.ClassName]](String message, [[.ErrorName]] error) {
		super(message);
		this.error = error;
	}

	public [[.ClassName]](String message, Throwable cause, [[.ErrorName]] error) {
		super(message, cause);
		this.error = error;
	}

	public [[.ClassName]](Throwable cause, [[.ErrorName]] error) {
		super(cause);
		this.error = error;
	}

	public [[.ErrorName]] getError() {
		return this.error;
	}
}
`)
	return w.ToCodeFile()
}
