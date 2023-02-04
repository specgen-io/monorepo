package test_service.models;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.*;
import org.junit.jupiter.api.Test;
import test_service.json.*;

import java.io.IOException;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;

import static org.junit.jupiter.api.Assertions.*;
import static test_service.utils.Utils.*;

public class JsonTest {
	private final Json json = createJson();

	public static Json createJson() {
		ObjectMapper objectMapper = new ObjectMapper();
		CustomObjectMapper.setup(objectMapper);
		return new Json(objectMapper);
	}

	public <T> void check(T data, String jsonStr, TypeReference<T> typeReference) {
		jsonStr = fixQuotes(jsonStr);

		var actualJson = json.write(data);
		assertEquals(jsonStr, actualJson);

		var actualData = json.read(jsonStr, typeReference);
		assertEquals(actualData, data);
	}

	@Test
	public void objectModel() {
		Message data = new Message(123);
		String jsonStr = "{'field':123}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void objectModelFieldCases() {
		MessageCases data = new MessageCases("snake_case value", "camelCase value");
		String jsonStr = "{'snake_case':'snake_case value','camelCase':'camelCase value'}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void objectModelMissingValueTypeField() {
		var exception = assertThrows(JsonParseException.class, () ->
			json.read("{}", new TypeReference<Message>() {})
		);
		assertNotNull(exception);
	}

	@Test
	public void nestedObject() {
		Parent data = new Parent("the string", new Message(123));
		String jsonStr = "{'field':'the string','nested':{'field':123}}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void objectFieldNotNull() {
		var exception = assertThrows(JsonParseException.class, () -> {
			var jsonStr = fixQuotes("{'field':'the string','nested':null}");
			json.read(jsonStr, new TypeReference<Parent>() {});
		});
		assertNotNull(exception);
	}

	@Test
	public void enumModel() {
		EnumFields data = new EnumFields(Choice.SECOND_CHOICE);
		String jsonStr = "{'enum_field':'Two'}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void numericTypes() {
		NumericFields data = new NumericFields(123, 1234, 1.23f, 1.23, new BigDecimal("1.23"));
		String jsonStr = "{'int_field':123,'long_field':1234,'float_field':1.23,'double_field':1.23,'decimal_field':1.23}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void nonNumericTypes() {
		NonNumericFields data = new NonNumericFields(true, "the string", UUID.fromString("123e4567-e89b-12d3-a456-426655440000"), LocalDate.parse("2019-11-30"), LocalDateTime.parse("2019-11-30T17:45:55"));
		String jsonStr = "{'boolean_field':true,'string_field':'the string','uuid_field':'123e4567-e89b-12d3-a456-426655440000','date_field':'2019-11-30','datetime_field':'2019-11-30T17:45:55'}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void arrayType() {
		ArrayFields data = new ArrayFields(Arrays.asList(1, 2, 3), (Arrays.asList("one", "two", "three")));
		String jsonStr = "{'int_array_field':[1,2,3],'string_array_field':['one','two','three']}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void mapType() {
    var map1 = new HashMap<String, Integer>();
    map1.put("one", 1);
    map1.put("two", 2);
    var map2 = new HashMap<String, String>();
    map2.put("one", "first");
    map2.put("two", "second");
		MapFields data = new MapFields(map1, map2);
		String jsonStr = "{'int_map_field':{'one':1,'two':2},'string_map_field':{'one':'first','two':'second'}}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void optionalTypes() {
		OptionalFields data = new OptionalFields(123, "the string");
		String jsonStr = "{'int_option_field':123,'string_option_field':'the string'}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void optionalTypesMissingFields() {
		var data = new OptionalFields(null, null);
		String jsonStr = "{}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void jsonType() throws IOException {
		String jsonField = fixQuotes("{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}");
		JsonNode node = new ObjectMapper().readTree(jsonField);
		RawJsonField data = new RawJsonField(node);
		String jsonStr = "{'json_field':{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void oneOfItemType() {
		OrderCreated data = new OrderCreated(UUID.fromString("58d5e212-165b-4ca0-909b-c86b9cee0111"), "SNI/01/136/0500", 3);
		String jsonStr = "{'id':'58d5e212-165b-4ca0-909b-c86b9cee0111','sku':'SNI/01/136/0500','quantity':3}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void oneOfWrapper() {
		OrderEventWrapper data = new OrderEventWrapper.Canceled(new OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")));
		String jsonStr = "{'canceled':{'id':'123e4567-e89b-12d3-a456-426655440000'}}";
		check(data, jsonStr, new TypeReference<>() {});
	}

	@Test
	public void oneOfItemNotNull() {
		assertThrows(JsonParseException.class, () -> {
			String jsonStr = "{'canceled':null}";
			json.read(jsonStr, new TypeReference<>() {});
		});
	}

	@Test
	public void jsonOneOfDiscriminatorTest() {
		OrderEventDiscriminator data = new OrderEventDiscriminator.Canceled(new OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")));
		String jsonStr = "{'_type':'canceled','id':'123e4567-e89b-12d3-a456-426655440000'}";
		check(data, jsonStr, new TypeReference<>() {});
	}
}
