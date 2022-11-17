package {{package.value}};

import org.springframework.context.annotation.*;
import springfox.documentation.builders.*;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger.web.*;

import java.util.*;

@Configuration
public class SpringFoxConfig {
	@Bean
	public Docket api() {
		return new Docket(DocumentationType.OAS_30)
			.select()
			.apis(RequestHandlerSelectors.any())
			.paths(PathSelectors.any())
			.build();
	}

	@Primary
	@Bean
	public SwaggerResourcesProvider swaggerResourcesProvider(InMemorySwaggerResourcesProvider defaultResourcesProvider) {
		return () -> {
			SwaggerResource wsResource = new SwaggerResource();
			wsResource.setName("documentation");
			wsResource.setSwaggerVersion("3.0");
			wsResource.setLocation("/docs/swagger.yaml");
			List<SwaggerResource> resources = new ArrayList<>();
			resources.add(wsResource);
			return resources;
		};
	}
}
