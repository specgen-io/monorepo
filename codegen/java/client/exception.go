package client

import (
	"fmt"
	"generator"
	"java/packages"
	"java/types"
	"java/writer"
	"spec"
)

func clientException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientException`)
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

func inheritedClientException(thePackage, errorsModelsPackage packages.Package, types *types.Types, error *spec.Response) *generator.CodeFile {
	errorName := types.Java(&error.Type.Definition)
	className := fmt.Sprintf(`%sException`, errorName)
	w := writer.New(thePackage, className)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: errorsModelsPackage.PackageName,
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
