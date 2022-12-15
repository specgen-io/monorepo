package test_client2.errors;

import test_client2.errors.models.*;

public class BadRequestErrorException extends ClientException {
	private final BadRequestError error;

	public BadRequestErrorException(BadRequestError error) {
		super("Body: %s" + error);
		this.error = error;
	}

	public BadRequestErrorException(String message, BadRequestError error) {
		super(message);
		this.error = error;
	}

	public BadRequestErrorException(String message, Throwable cause, BadRequestError error) {
		super(message, cause);
		this.error = error;
	}

	public BadRequestErrorException(Throwable cause, BadRequestError error) {
		super(cause);
		this.error = error;
	}

	public BadRequestError getError() {
		return this.error;
	}
}