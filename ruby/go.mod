module github.com/specgen-io/specgen/ruby/v2

go 1.18

replace github.com/specgen-io/specgen/yamlx/v2 => ../yamlx

replace github.com/specgen-io/specgen/spec/v2 => ../spec

replace github.com/specgen-io/specgen/generator/v2 => ../generator

require (
	github.com/pinzolo/casee v1.0.0
	github.com/specgen-io/specgen/generator/v2 v2.0.0-00010101000000-000000000000
	github.com/specgen-io/specgen/spec/v2 v2.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.5.0
)

require (
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/specgen-io/specgen/yamlx/v2 v2.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	gopkg.in/specgen-io/yaml.v3 v3.0.0-20211212030207-33c98a79c251 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
