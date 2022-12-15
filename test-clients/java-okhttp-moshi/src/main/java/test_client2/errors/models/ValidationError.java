package test_client2.errors.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class ValidationError {

	@Json(name = "path")
	private String path;

	@Json(name = "code")
	private String code;

	@Json(name = "message")
	private String message;

	public ValidationError(String path, String code, String message) {
		this.path = path;
		this.code = code;
		this.message = message;
	}

	public String getPath() {
		return path;
	}

	public void setPath(String path) {
		this.path = path;
	}

	public String getCode() {
		return code;
	}

	public void setCode(String code) {
		this.code = code;
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
		if (!(o instanceof ValidationError)) return false;
		ValidationError that = (ValidationError) o;
		return Objects.equals(getPath(), that.getPath()) && Objects.equals(getCode(), that.getCode()) && Objects.equals(getMessage(), that.getMessage());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getPath(), getCode(), getMessage());
	}

	@Override
	public String toString() {
		return String.format("ValidationError{path=%s, code=%s, message=%s}", path, code, message);
	}
}