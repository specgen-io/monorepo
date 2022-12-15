package test_client2.json;

import com.squareup.moshi.Moshi;
import test_client2.json.adapters.*;

public class CustomMoshiAdapters {
	public static Moshi.Builder setup(Moshi.Builder moshiBuilder) {
		moshiBuilder
			.add(new BigDecimalAdapter())
			.add(new UuidAdapter())
			.add(new LocalDateAdapter())
			.add(new LocalDateTimeAdapter());
		test_client2.errors.models.ModelsMoshiAdapters.setup(moshiBuilder);
		test_client2.v2.models.ModelsMoshiAdapters.setup(moshiBuilder);
		test_client2.models.ModelsMoshiAdapters.setup(moshiBuilder);

		return moshiBuilder;
	}
}