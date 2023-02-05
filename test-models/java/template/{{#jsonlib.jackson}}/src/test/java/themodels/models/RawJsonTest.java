package themodels.models;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.*;
import org.junit.jupiter.api.Test;
import themodels.json.*;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static themodels.utils.Utils.fixQuotes;

public class RawJsonTest {
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
    public void stringify() throws IOException {
        String jsonField = fixQuotes("{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}");
        JsonNode node = new ObjectMapper().readTree(jsonField);
        RawJsonField data = new RawJsonField(node);
        String expected = fixQuotes("RawJsonField{jsonField={'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}}");
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }

    @Test
    public void serialization() throws IOException {
        String jsonField = fixQuotes("{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}");
        JsonNode node = new ObjectMapper().readTree(jsonField);
        RawJsonField data = new RawJsonField(node);
        String jsonStr = "{'json_field':{'the_array':[1,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':123}}";
        check(data, jsonStr, new TypeReference<>() {});
    }
}
