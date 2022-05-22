package spec

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type operation struct {
	Endpoint     Endpoint     `yaml:"endpoint"`
	Description  *string      `yaml:"description,omitempty"`
	HeaderParams HeaderParams `yaml:"header,omitempty"`
	QueryParams  QueryParams  `yaml:"query,omitempty"`
	Body         *Definition  `yaml:"body,omitempty"`
	Responses    Responses    `yaml:"response"`
	Errors       Responses    `yaml:"errors"`
	Location     *yaml.Node
}

type Operation operation

func (operation *Operation) GetResponse(status string) *OperationResponse {
	for _, response := range operation.Responses {
		if response.Name.Source == status {
			return &response
		}
	}
	return nil
}

func (operation *Operation) BodyKind() BodyKind {
	return kindOf(operation.Body)
}

func (operation *Operation) BodyIs(kind BodyKind) bool {
	return kindOf(operation.Body) == kind
}

func (operation *Operation) HasParams() bool {
	return len(operation.QueryParams) > 0 || len(operation.HeaderParams) > 0 || len(operation.Endpoint.UrlParams) > 0
}

func (value *Operation) UnmarshalYAML(node *yaml.Node) error {
	internal := operation{}
	err := node.DecodeWith(decodeStrict, &internal)
	if err != nil {
		return err
	}
	internal.Location = node
	operation := Operation(internal)
	if operation.Body != nil && operation.Body.Description == nil {
		operation.Body.Description = getDescriptionFromComment(getMappingKey(node, "body"))
	}
	*value = operation
	return nil
}

func (value Operation) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("endpoint", value.Endpoint)
	yamlMap.AddOmitNil("description", value.Description)
	if len(value.HeaderParams) > 0 {
		yamlMap.Add("header", value.HeaderParams)
	}
	if len(value.QueryParams) > 0 {
		yamlMap.Add("query", value.QueryParams)
	}
	yamlMap.AddOmitNil("body", value.Body)
	yamlMap.Add("response", value.Responses)
	return yamlMap.Node, nil
}

type NamedOperation struct {
	Name Name
	Operation
	Api *Api
}

func (op *NamedOperation) FullUrl() string {
	return op.Api.Apis.GetUrl() + op.Endpoint.Url
}

func (op *NamedOperation) FullName() string {
	fullName := fmt.Sprintf(`%s.%s`, op.Api.Name.Source, op.Name.Source)
	if op.Api.Apis.Version.Version.Source != "" {
		fullName = fmt.Sprintf(`%s.%s`, op.Api.Apis.Version.Version.Source, fullName)
	}
	return fullName
}

type Operations []NamedOperation

func (value *Operations) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "operations should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedOperation, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		operation := Operation{}
		err = valueNode.DecodeWith(decodeStrict, &operation)
		if err != nil {
			return err
		}
		array[index] = NamedOperation{Name: name, Operation: operation}
	}
	*value = array
	return nil
}

func (value Operations) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(value); index++ {
		operation := value[index]
		yamlMap.Add(operation.Name, operation.Operation)
	}
	return yamlMap.Node, nil
}
