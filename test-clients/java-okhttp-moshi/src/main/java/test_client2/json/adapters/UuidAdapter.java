package test_client2.json.adapters;

import com.squareup.moshi.*;
import java.util.UUID;

public class UuidAdapter {
	@FromJson
	private UUID fromJson(String string) {
		return UUID.fromString(string);
	}

	@ToJson
	private String toJson(UUID value) {
		return value.toString();
	}
}