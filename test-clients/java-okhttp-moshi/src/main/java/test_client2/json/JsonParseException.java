package test_client2.json;

public class JsonParseException extends RuntimeException {
	public JsonParseException(Throwable exception) {
		super("Failed to parse body: " + exception.getMessage(), exception);
	}
}