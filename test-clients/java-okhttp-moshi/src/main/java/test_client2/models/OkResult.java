package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class OkResult {

	@Json(name = "ok_result")
	private String okResult;

	public OkResult(String okResult) {
		this.okResult = okResult;
	}

	public String getOkResult() {
		return okResult;
	}

	public void setOkResult(String okResult) {
		this.okResult = okResult;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof OkResult)) return false;
		OkResult that = (OkResult) o;
		return Objects.equals(getOkResult(), that.getOkResult());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getOkResult());
	}

	@Override
	public String toString() {
		return String.format("OkResult{okResult=%s}", okResult);
	}
}