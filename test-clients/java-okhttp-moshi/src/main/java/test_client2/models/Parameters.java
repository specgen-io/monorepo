package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class Parameters {

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

	@Json(name = "string_opt_field")
	private String stringOptField;

	@Json(name = "string_defaulted_field")
	private String stringDefaultedField;

	@Json(name = "string_array_field")
	private List<String> stringArrayField;

	@Json(name = "uuid_field")
	private UUID uuidField;

	@Json(name = "date_field")
	private LocalDate dateField;

	@Json(name = "date_array_field")
	private List<LocalDate> dateArrayField;

	@Json(name = "datetime_field")
	private LocalDateTime datetimeField;

	@Json(name = "enum_field")
	private Choice enumField;

	public Parameters(int intField, long longField, float floatField, double doubleField, BigDecimal decimalField, boolean boolField, String stringField, String stringOptField, String stringDefaultedField, List<String> stringArrayField, UUID uuidField, LocalDate dateField, List<LocalDate> dateArrayField, LocalDateTime datetimeField, Choice enumField) {
		this.intField = intField;
		this.longField = longField;
		this.floatField = floatField;
		this.doubleField = doubleField;
		this.decimalField = decimalField;
		this.boolField = boolField;
		this.stringField = stringField;
		this.stringOptField = stringOptField;
		this.stringDefaultedField = stringDefaultedField;
		this.stringArrayField = stringArrayField;
		this.uuidField = uuidField;
		this.dateField = dateField;
		this.dateArrayField = dateArrayField;
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

	public String getStringOptField() {
		return stringOptField;
	}

	public void setStringOptField(String stringOptField) {
		this.stringOptField = stringOptField;
	}

	public String getStringDefaultedField() {
		return stringDefaultedField;
	}

	public void setStringDefaultedField(String stringDefaultedField) {
		this.stringDefaultedField = stringDefaultedField;
	}

	public List<String> getStringArrayField() {
		return stringArrayField;
	}

	public void setStringArrayField(List<String> stringArrayField) {
		this.stringArrayField = stringArrayField;
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

	public List<LocalDate> getDateArrayField() {
		return dateArrayField;
	}

	public void setDateArrayField(List<LocalDate> dateArrayField) {
		this.dateArrayField = dateArrayField;
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
		if (!(o instanceof Parameters)) return false;
		Parameters that = (Parameters) o;
		return getIntField() == that.getIntField() && getLongField() == that.getLongField() && Float.compare(that.getFloatField(), getFloatField()) == 0 && Double.compare(that.getDoubleField(), getDoubleField()) == 0 && Objects.equals(getDecimalField(), that.getDecimalField()) && getBoolField() == that.getBoolField() && Objects.equals(getStringField(), that.getStringField()) && Objects.equals(getStringOptField(), that.getStringOptField()) && Objects.equals(getStringDefaultedField(), that.getStringDefaultedField()) && Objects.equals(getStringArrayField(), that.getStringArrayField()) && Objects.equals(getUuidField(), that.getUuidField()) && Objects.equals(getDateField(), that.getDateField()) && Objects.equals(getDateArrayField(), that.getDateArrayField()) && Objects.equals(getDatetimeField(), that.getDatetimeField()) && Objects.equals(getEnumField(), that.getEnumField());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getIntField(), getLongField(), getFloatField(), getDoubleField(), getDecimalField(), getBoolField(), getStringField(), getStringOptField(), getStringDefaultedField(), getStringArrayField(), getUuidField(), getDateField(), getDateArrayField(), getDatetimeField(), getEnumField());
	}

	@Override
	public String toString() {
		return String.format("Parameters{intField=%s, longField=%s, floatField=%s, doubleField=%s, decimalField=%s, boolField=%s, stringField=%s, stringOptField=%s, stringDefaultedField=%s, stringArrayField=%s, uuidField=%s, dateField=%s, dateArrayField=%s, datetimeField=%s, enumField=%s}", intField, longField, floatField, doubleField, decimalField, boolField, stringField, stringOptField, stringDefaultedField, stringArrayField, uuidField, dateField, dateArrayField, datetimeField, enumField);
	}
}