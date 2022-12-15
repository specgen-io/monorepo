package test_client2.errors.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class NotFoundError {

	@Json(name = "message")
	private String message;

	public NotFoundError(String message) {
		this.message = message;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof NotFoundError)) return false;
		NotFoundError that = (NotFoundError) o;
		return Objects.equals(getMessage(), that.getMessage());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getMessage());
	}

	@Override
	public String toString() {
		return String.format("NotFoundError{message=%s}", message);
	}
}