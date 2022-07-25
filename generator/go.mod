module github.com/specgen-io/specgen/generator/v2

go 1.18

replace github.com/specgen-io/specgen/yamlx/v2 => ../yamlx

replace github.com/specgen-io/specgen/spec/v2 => ../spec

require (
	github.com/specgen-io/specgen/spec/v2 v2.0.0-00010101000000-000000000000
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/pinzolo/casee v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/specgen-io/specgen/yamlx/v2 v2.0.0-00010101000000-000000000000 // indirect
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20211212030207-33c98a79c251 // indirect
)
