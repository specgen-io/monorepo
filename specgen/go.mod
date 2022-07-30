module github.com/specgen-io/specgen/v2

go 1.18

replace github.com/specgen-io/specgen/yamlx/v2 => ../yamlx

replace github.com/specgen-io/specgen/spec/v2 => ../spec

replace github.com/specgen-io/specgen/generator/v2 => ../generator

replace github.com/specgen-io/specgen/openapi/v2 => ../openapi

replace github.com/specgen-io/specgen/ruby/v2 => ../ruby

replace github.com/specgen-io/specgen/scala/v2 => ../scala

replace github.com/specgen-io/specgen/typescript/v2 => ../typescript

replace github.com/specgen-io/specgen/golang/v2 => ../golang

replace github.com/specgen-io/specgen/java/v2 => ../java

replace github.com/specgen-io/specgen/kotlin/v2 => ../kotlin

require (
	github.com/dollarshaveclub/line v0.0.0-20171219191008-fc7a351a8b58
	github.com/getkin/kin-openapi v0.85.0
	github.com/pinzolo/casee v1.0.0
	github.com/specgen-io/specgen/generator/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/golang/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/java/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/kotlin/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/openapi/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/ruby/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/scala/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/spec/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/typescript/v2 v2.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/specgen-io/specgen/yamlx/v2 v2.0.0-00010101000000-000000000000 // indirect
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20211212030207-33c98a79c251 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
