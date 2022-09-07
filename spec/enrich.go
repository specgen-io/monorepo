package spec

import (
	"errors"
	"fmt"
)

func enrich(specification *Spec) (*Messages, error) {
	messages := NewMessages()
	if specification.HttpErrors != nil {
		httpErrors := specification.HttpErrors
		httpErrors.InSpec = specification
		httpErrors.ResolvedModels = enrichModels(httpErrors.Models, messages)
		errorModels := buildModelsMap(httpErrors.Models)
		enricher := &httpEnricher{errorModels, messages}
		enricher.httpErrors(httpErrors)
	}
	for index := range specification.Versions {
		version := &specification.Versions[index]
		version.InSpec = specification
		version.ResolvedModels = enrichModels(version.Models, messages)
		models := buildModelsMap(version.Models)
		if specification.HttpErrors != nil {
			errorModels := buildModelsMap(specification.HttpErrors.Models)
			for name, model := range errorModels {
				models[name] = model
			}
		}
		enricher := &httpEnricher{models, messages}
		enricher.version(version)
	}
	if messages.ContainsLevel(LevelError) {
		return messages, errors.New("failed to parse specification")
	}
	return messages, nil
}

type httpEnricher struct {
	models   ModelsMap
	Messages *Messages
}

func (enricher *httpEnricher) version(version *Version) {
	for modIndex := range version.Models {
		version.Models[modIndex].InVersion = version
	}

	http := &version.Http
	http.InVersion = version

	for apiIndex := range http.Apis {
		api := &version.Http.Apis[apiIndex]
		api.InHttp = http
		for opIndex := range api.Operations {
			operation := &api.Operations[opIndex]
			operation.InApi = api
			enricher.operation(operation)
		}
	}
}

func (enricher *httpEnricher) httpErrors(httpErrors *HttpErrors) {
	for index := range httpErrors.Models {
		httpErrors.Models[index].InHttpErrors = httpErrors
	}

	for index := range httpErrors.Responses {
		enricher.definition(&httpErrors.Responses[index].Definition)
	}
}

func (enricher *httpEnricher) operation(operation *NamedOperation) {
	enricher.params(operation.Endpoint.UrlParams)
	enricher.params(operation.QueryParams)
	enricher.params(operation.HeaderParams)

	if operation.Body != nil {
		enricher.definition(operation.Body)
	}

	for index := range operation.Responses {
		operation.Responses[index].Operation = operation
		enricher.definition(&operation.Responses[index].Definition)
	}
}

func (enricher *httpEnricher) params(params []NamedParam) {
	for index := range params {
		enricher.typ(&params[index].DefinitionDefault.Type)
	}
}

func (enricher *httpEnricher) definition(definition *Definition) {
	if definition != nil {
		enricher.typ(&definition.Type)
	}
}

func (enricher *httpEnricher) typ(typ *Type) {
	enricher.typeDef(typ, &typ.Definition)
}

func (enricher *httpEnricher) typeDef(starter *Type, typ *TypeDef) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, found := enricher.models[typ.Plain]; found {
				typ.Info = ModelTypeInfo(model)
			} else {
				if info, found := Types[typ.Plain]; found {
					typ.Info = &info
				} else {
					e := Error("unknown type: %s", typ.Plain).At(locationFromNode(starter.Location))
					enricher.Messages.Add(e)
				}
			}
		case NullableType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = NullableTypeInfo(typ.Child.Info)
		case ArrayType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = ArrayTypeInfo()
		case MapType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = MapTypeInfo()
		default:
			panic(fmt.Sprintf("unknown kind of type: %v", typ))
		}
	}
}
