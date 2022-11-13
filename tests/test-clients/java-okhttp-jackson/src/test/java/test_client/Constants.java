package test_client;

import test_client.models.*;

import java.math.BigDecimal;
import java.time.*;
import java.util.*;

public class Constants {
	public static final String BASE_URL = "http://localhost:8081";

	public static final int INT_VALUE = 123;
	public static final long LONG_VALUE = 12345;
	public static final float FLOAT_VALUE = 1.23f;
	public static final double DOUBLE_VALUE = 12.345;
	public static final BigDecimal DECIMAL_VALUE = new BigDecimal("12345");
	public static final boolean BOOL_VALUE = true;
	public static final String STRING_VALUE = "the value";
	public static final String STRING_OPT_VALUE = "the value";
	public static final String STRING_DEFAULTED_VALUE = "value";
	public static final List<String> STRING_ARRAY_VALUE = Arrays.asList("the str1", "the str2");
	public static final UUID UUID_VALUE = UUID.fromString("123e4567-e89b-12d3-a456-426655440000");
	public static final LocalDate DATE_VALUE = LocalDate.parse("2020-01-01");
	public static final List<LocalDate> DATE_ARRAY_VALUE = Arrays.asList(LocalDate.parse("2020-01-01"), LocalDate.parse("2020-01-02"));
	public static final LocalDateTime DATETIME_VALUE = LocalDateTime.parse("2019-11-30T17:45:55");
	public static final Choice ENUM_VALUE = Choice.SECOND_CHOICE;

	public static final String BODY_STRING = "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu";
	public static final Message MESSAGE = new Message(INT_VALUE, STRING_VALUE);
	public static final HashMap<String, String> STRING_HASH_MAP = new HashMap<>() {{
		put("string_field", "the value");
		put("string_field_2", "the value_2");
	}};
}
