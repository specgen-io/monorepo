package {{package.value}};

import com.fasterxml.jackson.databind.ObjectMapper;
import io.micronaut.context.annotation.*;
import io.micronaut.jackson.ObjectMapperFactory;
import {{package.value}}.json.*;

@Factory
@Replaces(ObjectMapperFactory.class)
public class ObjectMapperConfig {
    @Bean
	@Replaces(ObjectMapper.class)
	public ObjectMapper objectMapper() {
		var objectMapper = new ObjectMapper();
		CustomObjectMapper.setup(objectMapper);
		return objectMapper;
    }

    @Bean
    public Json json() {
        return new Json(objectMapper());
    }
}
