package themodels.models;

import org.junit.jupiter.api.Test;
import themodels.models.RawJsonField;

import java.util.*;
import static themodels.models.Utils.*;

public class RawJsonTest extends JsonTest {
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
	  public void jsonType() {
		    //String jsonField = fixQuotes("{'the_array':[true,'some string'],'the_object':{'the_bool':true,'the_string':'some value'},'the_scalar':'the value'}");
		    var theObject = new HashMap<String, Object>();
        theObject.put("the_bool", true);
        theObject.put("the_string", "some value");
		    var theArray = new ArrayList<Object>(List.of(true, "some string"));
		    var map = new HashMap<String, Object>();
        map.put("the_array", theArray);
        map.put("the_object", theObject);
        map.put("the_scalar", "the value");
		    RawJsonField data = new RawJsonField(map);
		    String jsonStr = "{'json_field':{'the_array':[true,'some string'],'the_scalar':'the value','the_object':{'the_bool':true,'the_string':'some value'}}}";
		    check(data, jsonStr, RawJsonField.class);
	}
}