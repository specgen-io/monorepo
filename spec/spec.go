package spec

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
	"io/ioutil"
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

func specError(errs Messages) error {
	if len(errs) > 0 {
		message := ""
		for _, error := range errs {
			message = message + fmt.Sprintf("%s\n", error)
		}
		return errors.New("spec errors: \n" + message)
	}
	return nil
}

func ReadSpec(data []byte) (*SpecParseResult, error) {
	err := checkSpecVersion(data)
	if err != nil {
		return nil, err
	}

	spec, err := unmarshalSpec(data)
	if err != nil {
		return nil, err
	}

	errors := enrichSpec(spec)
	if len(errors) > 0 {
		return &SpecParseResult{Errors: errors}, specError(errors)
	}

	warnings, errors := validate(spec)
	if len(errors) > 0 {
		return &SpecParseResult{Spec: nil, Warnings: warnings, Errors: errors}, specError(errors)
	}

	return &SpecParseResult{Spec: spec, Warnings: warnings, Errors: errors}, nil
}

func ReadSpecFile(filepath string) (*SpecParseResult, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	result, err := ReadSpec(yamlFile)

	return result, err
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

func WriteSpecFile(spec *Spec, filepath string) error {
	data, err := WriteSpec(spec)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadMeta(data []byte) (*Meta, error) {
	err := checkSpecVersion(data)
	if err != nil {
		return nil, err
	}
	var meta Meta
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func ReadMetaFile(filepath string) (*Meta, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	meta, err := ReadMeta(yamlFile)

	if err != nil {
		return nil, err
	}

	return meta, nil
}
