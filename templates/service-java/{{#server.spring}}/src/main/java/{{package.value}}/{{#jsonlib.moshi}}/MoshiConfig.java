package {{package.value}};

import com.squareup.moshi.Moshi;
import org.springframework.context.annotation.*;
import {{package.value}}.json.*;

@Configuration
public class MoshiConfig {
	@Bean
	@Primary
	public Moshi getMoshi() {
		var moshiBuilder = new Moshi.Builder();
		CustomMoshiAdapters.setup(moshiBuilder);
		return moshiBuilder.build();
	}

    @Bean
    public Json json() {
        return new Json(getMoshi());
    }
}