package themodels.models;

import com.fasterxml.jackson.databind.*;
import org.junit.jupiter.api.Test;
import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static themodels.models.Utils.*;

public class RawJsonTest {
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
        check(data, jsonStr, RawJsonField.class);
    }
}
