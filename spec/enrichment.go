package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
)

func enrich(options SpecOptions, spec *Spec) (*Messages, error) {
	messages := NewMessages()
	for verIndex := range spec.Versions {
		version := &spec.Versions[verIndex]
		if options.AddErrors {
			if len(version.Http.Apis) > 0 {
				errorResponses, err := createErrorResponsesModels()
				if err != nil {
					return nil, err
				}
				version.Models = append(version.Models, errorResponses...)
				errors, err := createErrorResponses()
				if err != nil {
					return nil, err
				}
				version.Http.Errors = errors
			}
		}
		modelsMap := buildModelsMap(version.Models)
		enricher := &enricher{modelsMap, make(map[string]interface{}), messages, nil}
		enricher.Version(version)
		version.ResolvedModels = enricher.ResolvedModels
	}
	if messages.ContainsLevel(LevelError) {
		return messages, errors.New("failed to parse specification")
	}
	return messages, nil
}

type ModelsMap map[string]*NamedModel

func buildModelsMap(models Models) ModelsMap {
	result := make(map[string]*NamedModel)
	for modIndex := range models {
		name := models[modIndex].Name.Source
		result[name] = &models[modIndex]
	}
	return result
}

type enricher struct {
	ModelsMap       ModelsMap
	ProcessedModels map[string]interface{}
	Messages        *Messages
	ResolvedModels  []*NamedModel
}

func (enricher *enricher) processedModels(name string) bool {
	if _, ok := enricher.ProcessedModels[name]; ok {
		return true
	}
	enricher.ProcessedModels[name] = nil
	return false
}

func (enricher *enricher) findModel(name string) (*NamedModel, bool) {
	if model, ok := enricher.ModelsMap[name]; ok {
		return model, true
	}
	return nil, false
}

func (enricher *enricher) addResolvedModel(model *NamedModel) {
	for _, m := range enricher.ResolvedModels {
		if m.Name.Source == model.Name.Source {
			return
		}
	}
	enricher.ResolvedModels = append(enricher.ResolvedModels, model)
}

func (enricher *enricher) Version(version *Version) {
	for modIndex := range version.Models {
		model := &version.Models[modIndex]
		model.Version = version
		enricher.Model(model)
	}
	apis := &version.Http
	apis.Version = version

	for index := range apis.Errors {
		enricher.Definition(&apis.Errors[index].Definition)
	}

	for apiIndex := range apis.Apis {
		api := &version.Http.Apis[apiIndex]
		api.Http = apis
		for opIndex := range api.Operations {
			operation := &api.Operations[opIndex]
			operation.Api = api
			enricher.Operation(operation)
		}
	}
}

func (enricher *enricher) Operation(operation *NamedOperation) {
	enricher.Params(operation.Endpoint.UrlParams)
	enricher.Params(operation.QueryParams)
	enricher.Params(operation.HeaderParams)

	if operation.Body != nil {
		enricher.Definition(operation.Body)
	}

	for index := range operation.Responses {
		operation.Responses[index].Operation = operation
		enricher.Definition(&operation.Responses[index].Definition)
	}
}

func (enricher *enricher) Params(params []NamedParam) {
	for index := range params {
		enricher.DefinitionDefault(&params[index].DefinitionDefault)
	}
}

func (enricher *enricher) Model(model *NamedModel) {
	if !enricher.processedModels(model.Name.Source) {
		if model.IsObject() {
			for index := range model.Object.Fields {
				enricher.Definition(&model.Object.Fields[index].Definition)
			}
		}
		if model.IsOneOf() {
			for index := range model.OneOf.Items {
				enricher.Definition(&model.OneOf.Items[index].Definition)
			}
		}
		enricher.addResolvedModel(model)
	}
}

func (enricher *enricher) DefinitionDefault(definition *DefinitionDefault) {
	if definition != nil {
		enricher.Type(&definition.Type)
	}
}

func (enricher *enricher) Definition(definition *Definition) {
	if definition != nil {
		enricher.Type(&definition.Type)
	}
}

func (enricher *enricher) Type(typ *Type) {
	enricher.TypeDef(&typ.Definition, typ.Location)
}

func (enricher *enricher) TypeDef(typ *TypeDef, node *yaml.Node) *TypeInfo {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, ok := enricher.findModel(typ.Plain); ok {
				typ.Info = ModelTypeInfo(model)
				enricher.Model(model)
			} else {
				if info, ok := Types[typ.Plain]; ok {
					typ.Info = &info
				} else {
					e := Error("unknown type: %s", typ.Plain).At(locationFromNode(node))
					enricher.Messages.Add(e)
				}
			}
		case NullableType:
			childInfo := enricher.TypeDef(typ.Child, node)
			typ.Info = NullableTypeInfo(childInfo)
		case ArrayType:
			enricher.TypeDef(typ.Child, node)
			typ.Info = ArrayTypeInfo()
		case MapType:
			enricher.TypeDef(typ.Child, node)
			typ.Info = MapTypeInfo()
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
		return typ.Info
	}
	return nil
}

func createErrorResponsesModels() (Models, error) {
	data := `
BadRequestError:
  object:
    message: string
    location: ErrorLocation
    errors: ValidationError[]?

ValidationError:
  object:
    path: string
    code: string
    message: string?

ErrorLocation:
  enum:
    - query
    - header
    - body

NotFoundError:
  object:
    message: string

InternalServerError:
  object:
    message: string
`
	var models Models
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

const InternalServerError string = "InternalServerError"
const BadRequestError string = "BadRequestError"
const NotFoundError string = "NotFound"
const ValidationError string = "ValidationError"
const ErrorLocation string = "ErrorLocation"

func createErrorResponses() (Responses, error) {
	data := `
bad_request: BadRequestError   # Service will return this if parameters are not provided or couldn't be parsed correctly
not_found: NotFoundError   # Service will return this if the endpoint is not found
internal_server_error: InternalServerError   # Service will return this if unexpected internal error happens
`
	var responses Responses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	if err != nil {
		return nil, err
	}
	return responses, nil
}
