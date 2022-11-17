package theservice.services

{{#server.micronaut}}
import io.micronaut.context.annotation.Bean
{{/server.micronaut}}
{{#server.spring}}
import org.springframework.stereotype.Service
{{/server.spring}}
import theservice.models.*
import theservice.errors.models.*
import theservice.services.check.*

{{#server.micronaut}}
@Bean
{{/server.micronaut}}
{{#server.spring}}
@Service("CheckService")
{{/server.spring}}
class CheckServiceImpl : CheckService {
    override fun checkEmpty() {
        return
    }

    override fun checkEmptyResponse(body: Message) {
        return
    }

    override fun checkForbidden(): CheckForbiddenResponse {
        return CheckForbiddenResponse.Forbidden()
    }

    override fun sameOperationName(): SameOperationNameResponse {
        return SameOperationNameResponse.Ok()
    }

    override fun checkBadRequest(): CheckBadRequestResponse {
        val badRequestError = BadRequestError("Error returned from service implementation", ErrorLocation.UNKNOWN, null)
        return CheckBadRequestResponse.BadRequest(badRequestError)
    }
}
