package themodels.models;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.IOException;
import java.util.*;

import static themodels.models.Utils.*;
import themodels.models.RawJsonField;

public class RawJsonTest extends JsonTest {
    @Test
    public void stringify() throws IOException {
		    var theObject = new HashMap<String, Object>();
        theObject.put("the_bool", true);
        theObject.put("the_string", "some value");
		    var theArray = new ArrayList<Object>(List.of(1, "some string"));
		    var map = new HashMap<String, Object>();
        map.put("the_array", theArray);
        map.put("the_object", theObject);
        map.put("the_scalar", "the value");
        RawJsonField data = new RawJsonField(map);
        String expected = "RawJsonField{jsonField={the_array=[1, some string], the_scalar=the value, the_object={the_bool=true, the_string=some value}}}";
        String dataStr = data.toString();
        assertEquals(dataStr, expected);
    }
  
	  @Test
	  public void jsonType() {
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