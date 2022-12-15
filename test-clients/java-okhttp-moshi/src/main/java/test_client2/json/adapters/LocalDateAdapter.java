package test_client2.json.adapters;

import com.squareup.moshi.*;
import java.time.LocalDate;

public class LocalDateAdapter {
	@FromJson
	private LocalDate fromJson(String string) {
		return LocalDate.parse(string);
	}

	@ToJson
	private String toJson(LocalDate value) {
		return value.toString();
	}
}