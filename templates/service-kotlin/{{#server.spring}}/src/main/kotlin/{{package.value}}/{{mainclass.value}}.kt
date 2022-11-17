package {{package.value}}

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
{{#swagger.value}}
import springfox.documentation.swagger2.annotations.EnableSwagger2
{{/swagger.value}}

{{#swagger.value}}
@EnableSwagger2
{{/swagger.value}}
@SpringBootApplication
class {{mainclass.value}}

fun main(args: Array<String>) {
    runApplication<{{mainclass.value}}>(*args)
}
