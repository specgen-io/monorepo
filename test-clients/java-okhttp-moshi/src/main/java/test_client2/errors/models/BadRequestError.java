package test_client2.errors.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class BadRequestError {

	@Json(name = "message")
	private String message;

	@Json(name = "location")
	private ErrorLocation location;

	@Json(name = "errors")
	private List<ValidationError> errors;

	public BadRequestError(String message, ErrorLocation location, List<ValidationError> errors) {
		this.message = message;
		this.location = location;
		this.errors = errors;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	public ErrorLocation getLocation() {
		return location;
	}

	public void setLocation(ErrorLocation location) {
		this.location = location;
	}

	public List<ValidationError> getErrors() {
		return errors;
	}

	public void setErrors(List<ValidationError> errors) {
		this.errors = errors;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof BadRequestError)) return false;
		BadRequestError that = (BadRequestError) o;
		return Objects.equals(getMessage(), that.getMessage()) && Objects.equals(getLocation(), that.getLocation()) && Objects.equals(getErrors(), that.getErrors());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getMessage(), getLocation(), getErrors());
	}

	@Override
	public String toString() {
		return String.format("BadRequestError{message=%s, location=%s, errors=%s}", message, location, errors);
	}
}