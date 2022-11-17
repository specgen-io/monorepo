package {{package.value}}

import org.springframework.context.annotation.*
import springfox.documentation.builders.*
import springfox.documentation.spi.DocumentationType
import springfox.documentation.spring.web.plugins.Docket
import springfox.documentation.swagger.web.*

@Configuration
class SpringFoxConfig {
    @Bean
    fun api(): Docket {
        return Docket(DocumentationType.OAS_30)
            .select()
            .apis(RequestHandlerSelectors.any())
            .paths(PathSelectors.any())
            .build()
    }

    @Primary
    @Bean
    fun swaggerResourcesProvider(defaultResourcesProvider: InMemorySwaggerResourcesProvider): SwaggerResourcesProvider {
        return SwaggerResourcesProvider {
            val wsResource = SwaggerResource()
            wsResource.name = "documentation"
            wsResource.swaggerVersion = "3.0"
            wsResource.location = "/docs/swagger.yaml"
            val resources: MutableList<SwaggerResource> = ArrayList()
            resources.add(wsResource)
            resources
        }
    }
}
