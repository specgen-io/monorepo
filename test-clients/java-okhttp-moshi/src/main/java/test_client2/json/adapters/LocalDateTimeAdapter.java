package test_client2.json.adapters;

import com.squareup.moshi.*;
import java.time.LocalDateTime;

public class LocalDateTimeAdapter {
	@FromJson
	private LocalDateTime fromJson(String string) {
		return LocalDateTime.parse(string);
	}

	@ToJson
	private String toJson(LocalDateTime value) {
		return value.toString();
	}
}