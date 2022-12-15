package test_client2.errors;

import test_client2.errors.models.*;

public class InternalServerErrorException extends ClientException {
	private final InternalServerError error;

	public InternalServerErrorException(InternalServerError error) {
		super("Body: %s" + error);
		this.error = error;
	}

	public InternalServerErrorException(String message, InternalServerError error) {
		super(message);
		this.error = error;
	}

	public InternalServerErrorException(String message, Throwable cause, InternalServerError error) {
		super(message, cause);
		this.error = error;
	}

	public InternalServerErrorException(Throwable cause, InternalServerError error) {
		super(cause);
		this.error = error;
	}

	public InternalServerError getError() {
		return this.error;
	}
}