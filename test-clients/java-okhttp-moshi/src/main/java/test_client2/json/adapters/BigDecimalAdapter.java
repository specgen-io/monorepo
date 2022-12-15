package test_client2.json.adapters;

import com.squareup.moshi.*;
import okio.Okio;

import java.io.*;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;

public class BigDecimalAdapter {
	@FromJson
	public BigDecimal fromJson(JsonReader reader) throws IOException {
		var token = reader.peek();
		if (token != JsonReader.Token.NUMBER) {
			throw new JsonDataException("BigDecimal should be represented as number in JSON, found: "+token.name());
		}
		var source = reader.nextSource();
		return new BigDecimal(new String(source.readByteArray(), StandardCharsets.UTF_8));
	}

	@ToJson
	public void toJson(JsonWriter writer, BigDecimal value) throws IOException {
		var source = Okio.source(new ByteArrayInputStream(value.toString().getBytes()));
		var buffer = Okio.buffer(source);
		writer.value(buffer);
	}
}