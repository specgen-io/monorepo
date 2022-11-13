package theservice.services.v2;

import java.math.BigDecimal;
import java.time.*;
import java.util.*;

{{#server.micronaut}}
import io.micronaut.context.annotation.Bean;
{{/server.micronaut}}
{{#server.spring}}
import org.springframework.stereotype.Service;
{{/server.spring}}

import theservice.v2.models.*;
import theservice.v2.services.echo.*;

{{#server.micronaut}}
@Bean
{{/server.micronaut}}
{{#server.spring}}
@Service("EchoServiceV2")
{{/server.spring}}
public class EchoServiceImpl implements EchoService {
	@Override
	public Message echoBodyModel(Message body) {
		return body;
	}
}
