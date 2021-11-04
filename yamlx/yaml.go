package yamlx

import (
	"bytes"
	"gopkg.in/specgen-io/yaml.v3"
)

func ToYamlString(data interface{}) (string, error) {
	writer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		return "", nil
	}
	return writer.String(), nil
}
