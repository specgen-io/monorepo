module kotlin

go 1.18

replace yamlx => ../yamlx

replace spec => ../spec

replace generator => ../generator

replace openapi => ../openapi

require (
	github.com/pinzolo/casee v1.0.0
	generator v0.0.0-00010101000000-000000000000
	openapi v0.0.0-00010101000000-000000000000
	spec v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.5.0
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	yamlx v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20220807035601-846c18c37062 // indirect
)
