package spec

import (
	"bytes"
	"errors"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
	"strconv"
	"strings"
)

type Spec struct {
	Meta
	Versions []Version
}

type VersionSpecification struct {
	Http   Apis   `yaml:"http"`
	Models Models `yaml:"models"`
}

type Version struct {
	Version Name
	VersionSpecification
	ResolvedModels []*NamedModel
}

type Meta struct {
	SpecVersion string  `yaml:"spec"`
	Name        Name    `yaml:"name"`
	Title       *string `yaml:"title,omitempty"`
	Description *string `yaml:"description,omitempty"`
	Version     string  `yaml:"version"`
}

func isVersionNode(node *yaml.Node) bool {
	return VersionFormat.Check(node.Value) == nil
}

func (value *Spec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "http should be YAML mapping")
	}
	versions := []Version{}
	count := len(node.Content) / 2
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			version := Name{}
			err := keyNode.DecodeWith(decodeStrict, &version)
			if err != nil {
				return err
			}
			err = version.Check(VersionFormat)
			if err != nil {
				return err
			}

			theSpec := VersionSpecification{}
			valueNode.DecodeWith(decodeStrict, &theSpec)
			versions = append(versions, Version{version, theSpec, nil})
		}
	}
	theSpec := VersionSpecification{}
	err := node.DecodeWith(decodeLooze, &theSpec)
	if err != nil {
		return err
	}
	versions = append(versions, Version{Name{}, theSpec, nil})

	meta := Meta{}
	node.DecodeWith(decodeStrict, &meta)

	*value = Spec{meta, versions}
	return nil
}

func (value VersionSpecification) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	if len(value.Http.Apis) > 0 {
		yamlMap.Add("http", value.Http)
	}
	if len(value.Models) > 0 {
		yamlMap.Add("models", value.Models)
	}
	return yamlMap.Node, nil
}

func (value Spec) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Merge(value.Meta)
	for index := 0; index < len(value.Versions); index++ {
		version := value.Versions[index]

		if version.Version.Source != "" {
			yamlMap.Add(version.Version.Source, version.VersionSpecification)
		} else {
			yamlMap.Merge(version.VersionSpecification)
		}
	}
	return yamlMap.Node, nil
}

func (value Meta) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("spec", yamlx.String(value.SpecVersion))
	yamlMap.Add("name", value.Name)
	yamlMap.AddOmitNil("title", value.Title)
	yamlMap.AddOmitNil("description", value.Description)
	yamlMap.Add("version", yamlx.String(value.Version))
	return yamlMap.Node, nil
}

func unmarshalSpec(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.UnmarshalWith(decodeStrict, data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

type SpecParseResult struct {
	Spec     *Spec
	Warnings Messages
	Errors   Messages
}

func ReadSpec(data []byte) (*SpecParseResult, error) {
	data, err := checkSpecVersion(data)
	if err != nil {
		message := convertError(err)
		if message != nil {
			return &SpecParseResult{Errors: []Message{*message}}, errors.New("failed to read specification")
		}
		return nil, err
	}

	spec, err := unmarshalSpec(data)
	if err != nil {
		message := convertError(err)
		if message != nil {
			return &SpecParseResult{Errors: []Message{*message}}, errors.New("failed to read specification")
		}
		return nil, err
	}

	errs := enrichSpec(spec)
	if len(errs) > 0 {
		return &SpecParseResult{Errors: errs}, errors.New("failed to parse specification")
	}

	warns, errs := validate(spec)
	if len(errs) > 0 {
		return &SpecParseResult{Spec: nil, Warnings: warns, Errors: errs}, errors.New("failed to validate specification")
	}

	return &SpecParseResult{Spec: spec, Warnings: warns, Errors: errs}, nil
}

func convertError(err error) *Message {
	if strings.HasPrefix(err.Error(), "yaml: line ") {
		parts := strings.SplitN(strings.TrimPrefix(err.Error(), "yaml: line "), ":", 2)
		line, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil
		}
		return &Message{LevelError, strings.TrimSpace(parts[1]), line, 0}
	}
	return nil
}

func WriteSpec(spec *Spec) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(spec)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadMeta(data []byte) (*Meta, error) {
	data, err := checkSpecVersion(data)
	if err != nil {
		return nil, err
	}
	var meta Meta
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
