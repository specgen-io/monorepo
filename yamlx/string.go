package yamlx

import (
	"gopkg.in/specgen-io/yaml.v3"
)

func String(value string) yaml.Node {
	return yaml.Node{Kind: yaml.ScalarNode, Value: value}
}
