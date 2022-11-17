package {{package.value}};

import io.micronaut.runtime.Micronaut;
{{#swagger.value}}
import io.swagger.v3.oas.annotations.*;
import io.swagger.v3.oas.annotations.info.*;
{{/swagger.value}}

{{#swagger.value}}
@OpenAPIDefinition(
	info = @Info(
		title = "{{project.value}}",
		version = "1"
	)
)
{{/swagger.value}}
public class {{mainclass.value}} {

	public static void main(String[] args) {
		Micronaut.run({{mainclass.value}}.class, args);
	}
}
