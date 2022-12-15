package test_client2.errors;

import test_client2.errors.models.*;

public class NotFoundErrorException extends ClientException {
	private final NotFoundError error;

	public NotFoundErrorException(NotFoundError error) {
		super("Body: %s" + error);
		this.error = error;
	}

	public NotFoundErrorException(String message, NotFoundError error) {
		super(message);
		this.error = error;
	}

	public NotFoundErrorException(String message, Throwable cause, NotFoundError error) {
		super(message, cause);
		this.error = error;
	}

	public NotFoundErrorException(Throwable cause, NotFoundError error) {
		super(cause);
		this.error = error;
	}

	public NotFoundError getError() {
		return this.error;
	}
}