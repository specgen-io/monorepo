package {{package.value}};

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.context.annotation.*;
import {{package.value}}.json.*;

@Configuration
public class ObjectMapperConfig {
	@Bean
	@Primary
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
