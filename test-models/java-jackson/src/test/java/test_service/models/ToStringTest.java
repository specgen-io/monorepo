package test_service.models;

import com.fasterxml.jackson.databind.*;
import org.junit.jupiter.api.Test;

import java.io.*;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;

import static org.junit.jupiter.api.Assertions.*;
import static test_service.utils.Utils.*;

public class ToStringTest {

    @Test
    public void jsonMessageTest() {
        Message data = new Message(123);
        String expected = "Message{field=123}";
        assertEquals(data.toString(), expected);
    }

    @Test
    public void jsonParentTest() {
        Parent data = new Parent("the string", new Message(123));
        String expected = "Parent{field=the string, nested=Message{field=123}}";
        assertEquals(data.toString(), expected);
    }

    @Test
    public void jsonEnumTest() {
        EnumFields data = new EnumFields(Choice.SECOND_CHOICE);
        String expected = "EnumFields{enumField=SECOND_CHOICE}";
        assertEquals(data.toString(), expected);
    }

    @Test
    public void jsonNumericFieldsTest() {
        NumericFields data = new NumericFields(123, 1234, 1.23f, 1.23, new BigDecimal("1.23"));
        String expected = "NumericFields{intField=123, longField=1234, floatField=1.23, doubleField=1.23, decimalField=1.23}";
        assertEquals(data.toString(), expected);
    }

    @Test
    public void jsonNonNumericFieldsTest() {
        NonNumericFields data = new NonNumericFields(true, "the string", UUID.fromString("123e4567-e89b-12d3-a456-426655440000"), LocalDate.parse("2019-11-30"), LocalDateTime.parse("2019-11-30T17:45:55"));
        String expected = "NonNumericFields{booleanField=true, stringField=the string, uuidField=123e4567-e89b-12d3-a456-426655440000, dateField=2019-11-30, datetimeField=2019-11-30T17:45:55}";
        assertEquals(data.toString(), expected);
    }

    @Test
    public void jsonArrayFieldsTest() {
        ArrayFields data = new ArrayFields(Arrays.asList(1, 2, 3), (Arrays.asList("one", "two", "three")));
        String expected = "ArrayFields{intArrayField=[1, 2, 3], stringArrayField=[one, two, three]}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonMapFieldsTest() {
        MapFields data = new MapFields(new HashMap<String, Integer>() {{
            put("one", 1);
            put("two", 2);
        }}, new HashMap<String, String>() {{
            put("one", "first");
            put("two", "second");
        }});
        String expected = "MapFields{intMapField={one=1, two=2}, stringMapField={one=first, two=second}}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonOptionalFieldsTest() {
        OptionalFields data = new OptionalFields(123, "the string");
        String expected = "OptionalFields{intOptionField=123, stringOptionField=the string}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonRawJsonFieldTest() throws IOException {
        String jsonField = fixQuotes("{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}");
        JsonNode node = new ObjectMapper().readTree(jsonField);
        RawJsonField data = new RawJsonField(node);
        String expected = fixQuotes("RawJsonField{jsonField={'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}}");
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonOrderCreatedTest() {
        OrderCreated data = new OrderCreated(UUID.fromString("58d5e212-165b-4ca0-909b-c86b9cee0111"), "SNI/01/136/0500", 3);
        String expected = "OrderCreated{id=58d5e212-165b-4ca0-909b-c86b9cee0111, sku=SNI/01/136/0500, quantity=3}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonOneOfWrapperTest() {
        OrderEventWrapper data = new OrderEventWrapper.Canceled(new OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")));
        String expected = "Canceled{data=OrderCanceled{id=123e4567-e89b-12d3-a456-426655440000}}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void jsonOneOfDiscriminatorTest() {
        OrderEventDiscriminator data = new OrderEventDiscriminator.Canceled(new OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")));
        String expected = "Canceled{data=OrderCanceled{id=123e4567-e89b-12d3-a456-426655440000}}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }
}