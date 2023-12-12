package client

import (
	"fmt"
	"generator"
	"java/packages"
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

func responseException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ResponseException`)
	w.Lines(`
import java.lang.RuntimeException;

public class [[.ClassName]] extends RuntimeException {
	public [[.ClassName]](String message) {
		super(message);
	}
}
`)
	return w.ToCodeFile()
}

func errorResponseException(thePackage, errorsModelsPackage packages.Package, error *spec.Response) *generator.CodeFile {
	w := writer.New(thePackage, errorExceptionClassName(error))
	w.Imports.Star(errorsModelsPackage)
	w.Line(`public class [[.ClassName]] extends ResponseException {`)
	w.Indent()
	if !error.Body.IsEmpty() {
		errorBody := fmt.Sprintf(`%s body`, error.Body.Type.Definition)
		w.Line(`private final %s;`, errorBody)
		w.EmptyLine()
		w.Line(`public [[.ClassName]](%s) {`, errorBody)
		w.Line(`  super("Error response with status code %s");`, spec.HttpStatusCode(error.Name))
		w.Line(`  this.body = body;`)
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`public %s getBody() {`, error.Body.Type.Definition)
		w.Line(`	return this.body;`)
		w.Line(`}`)
	} else {
		w.Line(`public [[.ClassName]]() {`)
		w.Line(`  super("Error response with status code %s");`, spec.HttpStatusCode(error.Name))
		w.Line(`}`)
	}
	w.Unindent()
	w.Line(`}`)
	return w.ToCodeFile()
}

func errorExceptionClassName(error *spec.Response) string {
	return fmt.Sprintf(`%sException`, error.Name.PascalCase())
}
