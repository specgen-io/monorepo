package test_client2.models;

import com.squareup.moshi.Json;
import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;


public class AcceptedResult {

	@Json(name = "accepted_result")
	private String acceptedResult;

	public AcceptedResult(String acceptedResult) {
		this.acceptedResult = acceptedResult;
	}

	public String getAcceptedResult() {
		return acceptedResult;
	}

	public void setAcceptedResult(String acceptedResult) {
		this.acceptedResult = acceptedResult;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (!(o instanceof AcceptedResult)) return false;
		AcceptedResult that = (AcceptedResult) o;
		return Objects.equals(getAcceptedResult(), that.getAcceptedResult());
	}

	@Override
	public int hashCode() {
		return Objects.hash(getAcceptedResult());
	}

	@Override
	public String toString() {
		return String.format("AcceptedResult{acceptedResult=%s}", acceptedResult);
	}
}