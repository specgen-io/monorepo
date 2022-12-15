package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class UrlParameters {

	@Json(name = "int_field")
	private int intField;

	@Json(name = "long_field")
	private long longField;

	@Json(name = "float_field")
	private float floatField;

	@Json(name = "double_field")
	private double doubleField;

	@Json(name = "decimal_field")
	private BigDecimal decimalField;

	@Json(name = "bool_field")
	private boolean boolField;

	@Json(name = "string_field")
	private String stringField;

	@Json(name = "uuid_field")
	private UUID uuidField;

	@Json(name = "date_field")
	private LocalDate dateField;

	@Json(name = "datetime_field")
	private LocalDateTime datetimeField;

	@Json(name = "enum_field")
	private Choice enumField;

	public UrlParameters(int intField, long longField, float floatField, double doubleField, BigDecimal decimalField, boolean boolField, String stringField, UUID uuidField, LocalDate dateField, LocalDateTime datetimeField, Choice enumField) {
		this.intField = intField;
		this.longField = longField;
		this.floatField = floatField;
		this.doubleField = doubleField;
		this.decimalField = decimalField;
		this.boolField = boolField;
		this.stringField = stringField;
		this.uuidField = uuidField;
		this.dateField = dateField;
		this.datetimeField = datetimeField;
		this.enumField = enumField;
	}

	public int getIntField() {
		return intField;
	}

	public void setIntField(int intField) {
		this.intField = intField;
	}

	public long getLongField() {
		return longField;
	}

	public void setLongField(long longField) {
		this.longField = longField;
	}

	public float getFloatField() {
		return floatField;
	}

	public void setFloatField(float floatField) {
		this.floatField = floatField;
	}

	public double getDoubleField() {
		return doubleField;
	}

	public void setDoubleField(double doubleField) {
		this.doubleField = doubleField;
	}

	public BigDecimal getDecimalField() {
		return decimalField;
	}

	public void setDecimalField(BigDecimal decimalField) {
		this.decimalField = decimalField;
	}

	public boolean getBoolField() {
		return boolField;
	}

	public void setBoolField(boolean boolField) {
		this.boolField = boolField;
	}

	public String getStringField() {
		return stringField;
	}

	public void setStringField(String stringField) {
		this.stringField = stringField;
	}

	public UUID getUuidField() {
		return uuidField;
	}

	public void setUuidField(UUID uuidField) {
		this.uuidField = uuidField;
	}

	public LocalDate getDateField() {
		return dateField;
	}

	public void setDateField(LocalDate dateField) {
		this.dateField = dateField;
	}

	public LocalDateTime getDatetimeField() {
		return datetimeField;
	}

	public void setDatetimeField(LocalDateTime datetimeField) {
		this.datetimeField = datetimeField;
	}

	public Choice getEnumField() {
		return enumField;
	}

	public void setEnumField(Choice enumField) {
		this.enumField = enumField;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof UrlParameters)) return false;
		UrlParameters that = (UrlParameters) o;
		return getIntField() == that.getIntField() && getLongField() == that.getLongField() && Float.compare(that.getFloatField(), getFloatField()) == 0 && Double.compare(that.getDoubleField(), getDoubleField()) == 0 && Objects.equals(getDecimalField(), that.getDecimalField()) && getBoolField() == that.getBoolField() && Objects.equals(getStringField(), that.getStringField()) && Objects.equals(getUuidField(), that.getUuidField()) && Objects.equals(getDateField(), that.getDateField()) && Objects.equals(getDatetimeField(), that.getDatetimeField()) && Objects.equals(getEnumField(), that.getEnumField());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getIntField(), getLongField(), getFloatField(), getDoubleField(), getDecimalField(), getBoolField(), getStringField(), getUuidField(), getDateField(), getDatetimeField(), getEnumField());
	}

	@Override
	public String toString() {
		return String.format("UrlParameters{intField=%s, longField=%s, floatField=%s, doubleField=%s, decimalField=%s, boolField=%s, stringField=%s, uuidField=%s, dateField=%s, datetimeField=%s, enumField=%s}", intField, longField, floatField, doubleField, decimalField, boolField, stringField, uuidField, dateField, datetimeField, enumField);
	}
}