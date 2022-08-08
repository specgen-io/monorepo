module github.com/specgen-io/specgen/spec/v2

go 1.18

replace github.com/specgen-io/specgen/yamlx/v2 => ../yamlx

require (
	github.com/pinzolo/casee v1.0.0
	github.com/specgen-io/specgen/yamlx/v2 v2.0.0-00010101000000-000000000000
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20220807035601-846c18c37062
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
