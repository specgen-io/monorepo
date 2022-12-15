package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class Everything {

	@Json(name = "body_field")
	private Message bodyField;

	@Json(name = "float_query")
	private float floatQuery;

	@Json(name = "bool_query")
	private boolean boolQuery;

	@Json(name = "uuid_header")
	private UUID uuidHeader;

	@Json(name = "datetime_header")
	private LocalDateTime datetimeHeader;

	@Json(name = "date_url")
	private LocalDate dateUrl;

	@Json(name = "decimal_url")
	private BigDecimal decimalUrl;

	public Everything(Message bodyField, float floatQuery, boolean boolQuery, UUID uuidHeader, LocalDateTime datetimeHeader, LocalDate dateUrl, BigDecimal decimalUrl) {
		this.bodyField = bodyField;
		this.floatQuery = floatQuery;
		this.boolQuery = boolQuery;
		this.uuidHeader = uuidHeader;
		this.datetimeHeader = datetimeHeader;
		this.dateUrl = dateUrl;
		this.decimalUrl = decimalUrl;
	}

	public Message getBodyField() {
		return bodyField;
	}

	public void setBodyField(Message bodyField) {
		this.bodyField = bodyField;
	}

	public float getFloatQuery() {
		return floatQuery;
	}

	public void setFloatQuery(float floatQuery) {
		this.floatQuery = floatQuery;
	}

	public boolean getBoolQuery() {
		return boolQuery;
	}

	public void setBoolQuery(boolean boolQuery) {
		this.boolQuery = boolQuery;
	}

	public UUID getUuidHeader() {
		return uuidHeader;
	}

	public void setUuidHeader(UUID uuidHeader) {
		this.uuidHeader = uuidHeader;
	}

	public LocalDateTime getDatetimeHeader() {
		return datetimeHeader;
	}

	public void setDatetimeHeader(LocalDateTime datetimeHeader) {
		this.datetimeHeader = datetimeHeader;
	}

	public LocalDate getDateUrl() {
		return dateUrl;
	}

	public void setDateUrl(LocalDate dateUrl) {
		this.dateUrl = dateUrl;
	}

	public BigDecimal getDecimalUrl() {
		return decimalUrl;
	}

	public void setDecimalUrl(BigDecimal decimalUrl) {
		this.decimalUrl = decimalUrl;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof Everything)) return false;
		Everything that = (Everything) o;
		return Objects.equals(getBodyField(), that.getBodyField()) && Float.compare(that.getFloatQuery(), getFloatQuery()) == 0 && getBoolQuery() == that.getBoolQuery() && Objects.equals(getUuidHeader(), that.getUuidHeader()) && Objects.equals(getDatetimeHeader(), that.getDatetimeHeader()) && Objects.equals(getDateUrl(), that.getDateUrl()) && Objects.equals(getDecimalUrl(), that.getDecimalUrl());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getBodyField(), getFloatQuery(), getBoolQuery(), getUuidHeader(), getDatetimeHeader(), getDateUrl(), getDecimalUrl());
	}

	@Override
	public String toString() {
		return String.format("Everything{bodyField=%s, floatQuery=%s, boolQuery=%s, uuidHeader=%s, datetimeHeader=%s, dateUrl=%s, decimalUrl=%s}", bodyField, floatQuery, boolQuery, uuidHeader, datetimeHeader, dateUrl, decimalUrl);
	}
}