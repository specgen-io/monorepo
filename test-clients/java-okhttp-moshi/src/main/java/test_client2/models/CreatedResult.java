package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class CreatedResult {

	@Json(name = "created_result")
	private String createdResult;

	public CreatedResult(String createdResult) {
		this.createdResult = createdResult;
	}

	public String getCreatedResult() {
		return createdResult;
	}

	public void setCreatedResult(String createdResult) {
		this.createdResult = createdResult;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof CreatedResult)) return false;
		CreatedResult that = (CreatedResult) o;
		return Objects.equals(getCreatedResult(), that.getCreatedResult());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getCreatedResult());
	}

	@Override
	public String toString() {
		return String.format("CreatedResult{createdResult=%s}", createdResult);
	}
}