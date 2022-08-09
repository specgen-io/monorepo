module specgen/v2

go 1.18

replace yamlx => ../yamlx

replace spec => ../spec

replace generator => ../generator

replace openapi => ../openapi

replace ruby => ../ruby

replace scala => ../scala

replace typescript => ../typescript

replace golang => ../golang

replace java => ../java

replace kotlin => ../kotlin

require (
	generator v0.0.0-00010101000000-000000000000
	github.com/getkin/kin-openapi v0.85.0
	github.com/pinzolo/casee v1.0.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	golang v0.0.0-00010101000000-000000000000
	gotest.tools v2.2.0+incompatible
	java v0.0.0-00010101000000-000000000000
	kotlin v0.0.0-00010101000000-000000000000
	openapi v0.0.0-00010101000000-000000000000
	ruby v0.0.0-00010101000000-000000000000
	scala v0.0.0-00010101000000-000000000000
	spec v0.0.0-00010101000000-000000000000
	typescript v0.0.0-00010101000000-000000000000
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
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20220807035601-846c18c37062 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	yamlx v0.0.0-00010101000000-000000000000 // indirect
)
