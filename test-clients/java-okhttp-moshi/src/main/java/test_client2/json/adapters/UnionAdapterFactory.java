package test_client2.json.adapters;

import com.squareup.moshi.*;
import java.io.IOException;
import java.lang.annotation.Annotation;
import java.lang.reflect.Type;
import java.util.*;
import javax.annotation.*;

public final class UnionAdapterFactory<T> implements JsonAdapter.Factory {
		final Class<T> baseType;
		final String discriminator;
		final List<String> tags;
		final List<Type> subtypes;
		@Nullable final JsonAdapter<Object> fallbackAdapter;

		UnionAdapterFactory(
						Class<T> baseType,
						String discriminator,
						List<String> tags,
						List<Type> subtypes,
						@Nullable JsonAdapter<Object> fallbackAdapter) {
				this.baseType = baseType;
				this.discriminator = discriminator;
				this.tags = tags;
				this.subtypes = subtypes;
				this.fallbackAdapter = fallbackAdapter;
		}

		/**
		 * @param baseType The base type for which this factory will create adapters. Cannot be Object.
		 *     JSON object.
		 */
		@CheckReturnValue
		public static <T> UnionAdapterFactory<T> of(Class<T> baseType) {
				if (baseType == null) throw new NullPointerException("baseType == null");
				return new UnionAdapterFactory<>(baseType, null, Collections.<String>emptyList(), Collections.<Type>emptyList(), null);
		}

		/** Returns a new factory that decodes instances of {@code subtype}. */
		public UnionAdapterFactory<T> withDiscriminator(String discriminator) {
				if (discriminator == null) throw new NullPointerException("discriminator == null");
				return new UnionAdapterFactory<>(baseType, discriminator, tags, subtypes, fallbackAdapter);
		}

		/** Returns a new factory that decodes instances of {@code subtype}. */
		public UnionAdapterFactory<T> withSubtype(Class<? extends T> subtype, String tag) {
				if (subtype == null) throw new NullPointerException("subtype == null");
				if (tag == null) throw new NullPointerException("tag == null");
				if (tags.contains(tag)) {
						throw new IllegalArgumentException("Tags must be unique.");
				}
				List<String> newTags = new ArrayList<>(tags);
				newTags.add(tag);
				List<Type> newSubtypes = new ArrayList<>(subtypes);
				newSubtypes.add(subtype);
				return new UnionAdapterFactory<>(baseType, discriminator, newTags, newSubtypes, fallbackAdapter);
		}

		/**
		 * Returns a new factory that with default to {@code fallbackJsonAdapter.fromJson(reader)} upon
		 * decoding of unrecognized tags.
		 *
		 * <p>The {@link JsonReader} instance will not be automatically consumed, so make sure to consume
		 * it within your implementation of {@link JsonAdapter#fromJson(JsonReader)}
		 */
		public UnionAdapterFactory<T> withFallbackAdapter(@Nullable JsonAdapter<Object> fallbackJsonAdapter) {
				return new UnionAdapterFactory<>(baseType, discriminator, tags, subtypes, fallbackJsonAdapter);
		}

		/**
		 * Returns a new factory that will default to {@code defaultValue} upon decoding of unrecognized
		 * tags. The default value should be immutable.
		 */
		public UnionAdapterFactory<T> withDefaultValue(@Nullable T defaultValue) {
				return withFallbackAdapter(buildFallbackAdapter(defaultValue));
		}

		private JsonAdapter<Object> buildFallbackAdapter(final T defaultValue) {
				return new JsonAdapter<Object>() {
						@Override
						public @Nullable Object fromJson(JsonReader reader) throws IOException {
								reader.skipValue();
								return defaultValue;
						}

						@Override
						public void toJson(JsonWriter writer, Object value) throws IOException {
								throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
						}
				};
		}

		@Override
		public JsonAdapter<?> create(Type type, Set<? extends Annotation> annotations, Moshi moshi) {
				if (Types.getRawType(type) != baseType || !annotations.isEmpty()) {
						return null;
				}

				List<JsonAdapter<Object>> jsonAdapters = new ArrayList<>(subtypes.size());
				for (int i = 0, size = subtypes.size(); i < size; i++) {
						jsonAdapters.add(moshi.adapter(subtypes.get(i)));
				}

				if (discriminator != null) {
						return new UnionDiscriminatorAdapter(discriminator, tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe();
				} else {
						return new UnionWrapperAdapter(tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe();
				}
		}

		static final class UnionDiscriminatorAdapter extends JsonAdapter<Object> {
				final String discriminator;
				final List<String> tags;
				final List<Type> subtypes;
				final List<JsonAdapter<Object>> adapters;
				@Nullable final JsonAdapter<Object> fallbackAdapter;
				final JsonReader.Options discriminatorOptions;
				final JsonReader.Options tagsOptions;

				UnionDiscriminatorAdapter(
								String discriminator,
								List<String> tags,
								List<Type> subtypes,
								List<JsonAdapter<Object>> adapters,
								@Nullable JsonAdapter<Object> fallbackAdapter) {
						this.discriminator = discriminator;
						this.tags = tags;
						this.subtypes = subtypes;
						this.adapters = adapters;
						this.fallbackAdapter = fallbackAdapter;

						this.discriminatorOptions = JsonReader.Options.of(discriminator);
						this.tagsOptions = JsonReader.Options.of(tags.toArray(new String[0]));
				}

				@Override
				public Object fromJson(JsonReader reader) throws IOException {
						final int tagIndex = getTagIndex(reader);
						JsonAdapter<Object> adapter = fallbackAdapter;
						if (tagIndex != -1) {
								adapter = adapters.get(tagIndex);
						}
						return adapter.fromJson(reader);
				}

				@Override
				public void toJson(JsonWriter writer, Object value) throws IOException {
						int tagIndex = getTagIndex(value);
						if (tagIndex == -1) {
								fallbackAdapter.toJson(writer, value);
						} else {
								final JsonAdapter<Object> adapter = adapters.get(tagIndex);
								writer.beginObject();
								writer.name(discriminator).value(tags.get(tagIndex));
								int flattenToken = writer.beginFlatten();
								adapter.toJson(writer, value);
								writer.endFlatten(flattenToken);
								writer.endObject();
						}
				}

				private int getTagIndex(JsonReader reader) throws IOException {
						JsonReader peeked = reader.peekJson();
						peeked.setFailOnUnknown(false);
						try {
								peeked.beginObject();
								while (peeked.hasNext()) {
										if (peeked.selectName(discriminatorOptions) == -1) {
												peeked.skipName();
												peeked.skipValue();
												continue;
										}

										int tagIndex = peeked.selectString(tagsOptions);
										if (tagIndex == -1 && this.fallbackAdapter == null) {
												throw new JsonDataException("Expected one of " + tags + " for key '" + discriminator + "' but found '" + peeked.nextString() + "'. Register a subtype for this tag.");
										}
										return tagIndex;
								}

								throw new JsonDataException("Missing discriminator field " + discriminator);
						} finally {
								peeked.close();
						}
				}

				private int getTagIndex(Object value) {
						Class<?> type = value.getClass();
						int tagIndex = subtypes.indexOf(type);
						if (tagIndex == -1 && fallbackAdapter == null) {
								throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
						} else {
								return tagIndex;
						}
				}

				@Override
				public String toString() {
						return "UnionDiscriminatorAdapter(" + discriminator + ")";
				}
		}

		static final class UnionWrapperAdapter extends JsonAdapter<Object> {
				final List<String> tags;
				final List<Type> subtypes;
				final List<JsonAdapter<Object>> adapters;
				@Nullable final JsonAdapter<Object> fallbackAdapter;
				final JsonReader.Options tagsOptions;

				UnionWrapperAdapter(
								List<String> tags,
								List<Type> subtypes,
								List<JsonAdapter<Object>> adapters,
								@Nullable JsonAdapter<Object> fallbackAdapter) {
						this.tags = tags;
						this.subtypes = subtypes;
						this.adapters = adapters;
						this.fallbackAdapter = fallbackAdapter;

						this.tagsOptions = JsonReader.Options.of(tags.toArray(new String[0]));
				}

				@Override
				public Object fromJson(JsonReader reader) throws IOException {
						int tagIndex = getTagIndex(reader);
						if (tagIndex == -1) {
								return this.fallbackAdapter.fromJson(reader);
						} else {
								reader.beginObject();
								reader.skipName();
								final Object value = adapters.get(tagIndex).fromJson(reader);
								reader.endObject();
								return value;
						}
				}

				@Override
				public void toJson(JsonWriter writer, Object value) throws IOException {
						int tagIndex = getTagIndex(value);
						if (tagIndex == -1) {
								fallbackAdapter.toJson(writer, value);
						} else {
								final JsonAdapter<Object> adapter = adapters.get(tagIndex);
								writer.beginObject();
								writer.name(tags.get(tagIndex));
								adapter.toJson(writer, value);
								writer.endObject();
						}
				}

				private int getTagIndex(JsonReader reader) throws IOException {
						JsonReader peeked = reader.peekJson();
						peeked.setFailOnUnknown(false);
						try {
								peeked.beginObject();
								int tagIndex = peeked.selectName(tagsOptions);
								if (tagIndex == -1 && this.fallbackAdapter == null) {
										throw new JsonDataException("Expected one of keys:" + tags + "' but found '" + peeked.nextString() + "'. Register a subtype for this tag.");
								}
								return tagIndex;
						} finally {
								peeked.close();
						}
				}

				private int getTagIndex(Object value) {
						Class<?> type = value.getClass();
						int tagIndex = subtypes.indexOf(type);
						if (tagIndex == -1 && fallbackAdapter == null) {
								throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
						}
						return tagIndex;
				}

				@Override
				public String toString() {
						return "UnionWrapperAdapter";
				}
		}
}