package test_client2.errors;

import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.lang.reflect.ParameterizedType;
import java.io.IOException;
import okhttp3.*;
import org.jetbrains.annotations.NotNull;
import org.slf4j.*;
import test_client2.errors.models.*;
import test_client2.json.*;
import static test_client2.utils.ClientResponse.*;

public class ErrorsInterceptor implements Interceptor {
	private static final Logger logger = LoggerFactory.getLogger(ErrorsInterceptor.class);

	private final Json json;

	public ErrorsInterceptor(Json json) {
		this.json = json;
	}
	
	@NotNull
	@Override
	public Response intercept(@NotNull Chain chain) throws IOException {
		var request = chain.request();
		var response = chain.proceed(request);
		if (response.code() == 400) {
			var responseBodyString = getResponseBodyString(response, logger);
			var responseBody = json.read(responseBodyString, BadRequestError.class);
			throw new BadRequestErrorException(responseBody);
		}
		if (response.code() == 404) {
			var responseBodyString = getResponseBodyString(response, logger);
			var responseBody = json.read(responseBodyString, NotFoundError.class);
			throw new NotFoundErrorException(responseBody);
		}
		if (response.code() == 500) {
			var responseBodyString = getResponseBodyString(response, logger);
			var responseBody = json.read(responseBodyString, InternalServerError.class);
			throw new InternalServerErrorException(responseBody);
		}
		return response;
	}
}