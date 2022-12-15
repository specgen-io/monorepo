package test_client2.json;

import com.squareup.moshi.Moshi;

import java.lang.reflect.ParameterizedType;

public class Json {
	private final Moshi moshi;

	public Json(Moshi moshi) {
		this.moshi = moshi;
	}

	public <T> String write(Class<T> type, T data) {
		return moshi.adapter(type).toJson(data);
	}

	public <T> String write(ParameterizedType type, T data) {
		return moshi.adapter(type).toJson(data);
	}

	public <T> T read(String jsonStr, Class<T> type) {
		try {
			return moshi.adapter(type).fromJson(jsonStr);
		} catch (Exception exception) {
			throw new JsonParseException(exception);
		}
	}

	public <T> T read(String jsonStr, ParameterizedType type) {
		try {
			return moshi.<T>adapter(type).fromJson(jsonStr);
		} catch (Exception exception) {
			throw new JsonParseException(exception);
		}
	}
}