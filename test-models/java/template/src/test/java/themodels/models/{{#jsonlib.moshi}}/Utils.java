package themodels.models;

import com.squareup.moshi.*;
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
}