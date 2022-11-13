package theservice.services;

import java.math.BigDecimal;
import java.time.*;
import java.util.*;

{{#server.micronaut}}
import io.micronaut.context.annotation.Bean;
{{/server.micronaut}}
{{#server.spring}}
import org.springframework.stereotype.Service;
{{/server.spring}}

import theservice.models.*;
import theservice.errors.models.*;
import theservice.services.check.*;

{{#server.micronaut}}
@Bean
{{/server.micronaut}}
{{#server.spring}}
@Service("CheckService")
{{/server.spring}}
public class CheckServiceImpl implements CheckService {
	@Override
	public void checkEmpty() {
		return;
	}

	@Override
	public void checkEmptyResponse(Message body) {
		return;
	}
	
	@Override
	public CheckForbiddenResponse checkForbidden() {
		return new CheckForbiddenResponse.Forbidden();
	}

	@Override
	public SameOperationNameResponse sameOperationName() {
		return new SameOperationNameResponse.Ok();
	}

  @Override
  public CheckBadRequestResponse checkBadRequest() {
    var badRequestError = new BadRequestError("Error returned from service implementation", ErrorLocation.UNKNOWN, null);
    return new CheckBadRequestResponse.BadRequest(badRequestError);
  }
}
