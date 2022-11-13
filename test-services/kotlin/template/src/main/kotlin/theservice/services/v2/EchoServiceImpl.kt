package theservice.services.v2

{{#server.micronaut}}
import io.micronaut.context.annotation.Bean
{{/server.micronaut}}
{{#server.spring}}
import org.springframework.stereotype.Service
{{/server.spring}}
import theservice.v2.models.*
import theservice.v2.services.echo.*
import java.math.BigDecimal
import java.time.*
import java.util.*
import java.io.*

{{#server.micronaut}}
@Bean
{{/server.micronaut}}
{{#server.spring}}
@Service("EchoServiceV2")
{{/server.spring}}
class EchoServiceImpl : EchoService {
	override fun echoBodyModel(body: Message): Message {
        return body
	}
}
