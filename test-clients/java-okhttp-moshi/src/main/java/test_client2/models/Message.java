package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class Message {

	@Json(name = "int_field")
	private int intField;

	@Json(name = "string_field")
	private String stringField;

	public Message(int intField, String stringField) {
		this.intField = intField;
		this.stringField = stringField;
	}

	public int getIntField() {
		return intField;
	}

	public void setIntField(int intField) {
		this.intField = intField;
	}

	public String getStringField() {
		return stringField;
	}

	public void setStringField(String stringField) {
		this.stringField = stringField;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof Message)) return false;
		Message that = (Message) o;
		return getIntField() == that.getIntField() && Objects.equals(getStringField(), that.getStringField());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getIntField(), getStringField());
	}

	@Override
	public String toString() {
		return String.format("Message{intField=%s, stringField=%s}", intField, stringField);
	}
}