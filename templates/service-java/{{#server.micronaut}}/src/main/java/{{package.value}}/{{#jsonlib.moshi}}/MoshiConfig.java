package {{package.value}};

import com.squareup.moshi.Moshi;
import io.micronaut.context.annotation.*;
import {{package.value}}.json.*;

@Factory
public class MoshiConfig {
    @Bean
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