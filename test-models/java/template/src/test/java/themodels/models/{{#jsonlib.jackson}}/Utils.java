package themodels.models;

import com.fasterxml.jackson.databind.*;
import org.json.JSONException;
import org.skyscreamer.jsonassert.*;
import themodels.json.*;
import themodels.json.Json;

import static org.junit.jupiter.api.Assertions.*;

public class Utils {
    public static String fixQuotes(String jsonStr) {
        return jsonStr.replaceAll("'", "\"");
    }

    public static Json createJson() {
        ObjectMapper objectMapper = new ObjectMapper();
        CustomObjectMapper.setup(objectMapper);
        return new Json(objectMapper);
    }

    public static <T> void check(T data, String jsonStr, Class<T> tClass) {
        var json = createJson();
        jsonStr = fixQuotes(jsonStr);

        var actualJson = json.write(data);

        try {
            JSONAssert.assertEquals(jsonStr, actualJson, JSONCompareMode.NON_EXTENSIBLE);
        } catch (JSONException ex) {
            fail(ex);
        }

        var actualData = json.read(jsonStr, tClass);
        assertEquals(actualData, data);
    }
}