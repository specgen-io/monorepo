package themodels.models;

{{#jsonlib.jackson}}
import com.fasterxml.jackson.databind.*;
{{/jsonlib.jackson}}
{{#jsonlib.moshi}}
import com.squareup.moshi.*;
{{/jsonlib.moshi}}
import org.json.JSONException;
import org.skyscreamer.jsonassert.*;
import themodels.json.*;

import static org.junit.jupiter.api.Assertions.*;

public class Utils {
    public static String fixQuotes(String jsonStr) {
        return jsonStr.replaceAll("'", "\"");
    }

    {{#jsonlib.jackson}}
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
    {{/jsonlib.jackson}}
    {{#jsonlib.moshi}}
    public static Json createJson() {
        Moshi.Builder moshiBuilder = new Moshi.Builder();
        CustomMoshiAdapters.setup(moshiBuilder);
        return new Json(moshiBuilder.build());
    }
    
    public static <T> void check(T data, String jsonStr, Class<T> tClass) {
        var json = createJson();
        jsonStr = fixQuotes(jsonStr);

        var actualJson = json.write(tClass, data);
        try {
            JSONAssert.assertEquals(jsonStr, actualJson, JSONCompareMode.NON_EXTENSIBLE);
        } catch (JSONException ex) {
            fail(ex);
        }

        var actualData = json.read(jsonStr, tClass);
        assertEquals(actualData, data);
    }
    {{/jsonlib.moshi}}
}