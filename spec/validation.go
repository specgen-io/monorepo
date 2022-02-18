package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"strconv"
	"strings"
)

func (self Message) String() string {
	if self.Location != nil {
		return fmt.Sprintf("%s - line %s: %s", self.Level, strconv.Itoa(self.Location.Line), self.Message)
	}
	return fmt.Sprintf("%s - %s", self.Level, self.Message)
}

func validate(spec *Spec) (*Messages, error) {
	messages := NewMessages()
	validator := &validator{messages}
	validator.Spec(spec)
	var err error = nil
	if validator.Messages.ContainsLevel(LevelError) {
		err = errors.New("failed to validate specification")
	}
	return messages, err
}

type validator struct {
	Messages *Messages
}

func (validator *validator) addError(node *yaml.Node, message string) {
	validator.Messages.Add(Error(message).At(locationFromNode(node)))
}

func (validator *validator) addWarning(node *yaml.Node, message string) {
	validator.Messages.Add(Warning(message).At(locationFromNode(node)))
}

func (validator *validator) Spec(spec *Spec) {
	urlsMap := map[string][]*NamedOperation{}
	for versionIndex := range spec.Versions {
		models := spec.Versions[versionIndex].Models
		for modelIndex := range models {
			validator.Model(&models[modelIndex])
		}
		apis := spec.Versions[versionIndex].Http.Apis
		for apiIndex := range apis {
			for operationIndex := range apis[apiIndex].Operations {
				operation := apis[apiIndex].Operations[operationIndex]
				url := operation.Endpoint.Method + " " + operation.FullUrl()
				if _, found := urlsMap[url]; !found {
					urlsMap[url] = []*NamedOperation{}
				}
				urlsMap[url] = append(urlsMap[url], &operation)
				validator.Operation(&operation)
			}
		}
	}
	for url, operations := range urlsMap {
		if len(operations) > 1 {
			operationsStr := []string{}
			for _, operation := range operations {
				operationsStr = append(operationsStr, operation.FullName())
			}
			validator.addWarning(operations[0].Location, fmt.Sprintf(`endpoint "%s" is used for %d operations: %s`, url, len(operationsStr), strings.Join(operationsStr, ", ")))
		}
	}
}

func (validator *validator) ParamsNames(paramsMap map[string]NamedParam, params []NamedParam) {
	for _, p := range params {
		if other, ok := paramsMap[p.Name.SnakeCase()]; ok {
			message := fmt.Sprintf("parameter name '%s' conflicts with the other parameter name '%s'", p.Name.Source, other.Name.Source)
			validator.addError(p.Name.Location, message)
		} else {
			paramsMap[p.Name.SnakeCase()] = p
		}
	}
}

func (validator *validator) Operation(operation *NamedOperation) {
	paramsMap := make(map[string]NamedParam)
	validator.ParamsNames(paramsMap, operation.Endpoint.UrlParams)
	validator.ParamsNames(paramsMap, operation.QueryParams)
	validator.ParamsNames(paramsMap, operation.HeaderParams)

	validator.Params(operation.Endpoint.UrlParams, false)
	validator.Params(operation.QueryParams, true)
	validator.Params(operation.HeaderParams, true)

	if operation.Body != nil && !operation.Body.Type.Definition.IsEmpty() {
		bodyType := operation.Body.Type
		if bodyType.Definition.Info.Structure != StructureObject &&
			bodyType.Definition.Info.Structure != StructureArray &&
			bodyType.Definition.String() != TypeString {
			message := fmt.Sprintf("body should be object, array or string type, found %s", bodyType.Definition.Name)
			validator.addError(operation.Body.Location, message)
		}
		validator.Definition(operation.Body)
	}

	for index := range operation.Responses {
		validator.Response(&operation.Responses[index])
	}
}

func (validator *validator) Response(response *NamedResponse) {
	if response.Name.Source == HttpStatusInternalServerError || response.Name.Source == HttpStatusNotFound || response.Name.Source == HttpStatusBadRequest {
		if !response.Type.Definition.IsEmpty() {
			message := fmt.Sprintf("response %s can be only empty if declared, found %s", response.Name.Source, response.Type.Definition.Name)
			validator.addWarning(response.Type.Location, message)
		}
	}
	if !response.Type.Definition.IsEmpty() &&
		response.Type.Definition.Info.Structure != StructureObject &&
		response.Type.Definition.Info.Structure != StructureArray &&
		response.Type.Definition.String() != TypeString {
		message := fmt.Sprintf("response %s should be either empty or some type with structure of an object or array, found %s", response.Name.Source, response.Type.Definition.Name)
		validator.addError(response.Type.Location, message)
	}
	validator.Definition(&response.Definition)
}

func (validator *validator) Params(params []NamedParam, allowArrayTypes bool) {
	for index := range params {
		paramName := params[index].Name
		paramType := params[index].DefinitionDefault.Type
		scalar := paramType.Definition.Info.Structure == StructureScalar
		arrayNotNullable := paramType.Definition.Info.Structure == StructureArray && !paramType.Definition.IsNullable()
		if allowArrayTypes {
			if !scalar && !arrayNotNullable {
				validator.addError(paramType.Location, fmt.Sprintf("parameter %s should be of scalar type or array of scalar type, found %s", paramName.Source, paramType.Definition.Name))
			}
		} else {
			if !scalar {
				validator.addError(paramType.Location, fmt.Sprintf("parameter %s should be of scalar type, found %s", paramName.Source, paramType.Definition.Name))
			}
		}
		validator.DefinitionDefault(&params[index].DefinitionDefault)
	}
}

func (validator *validator) ItemsUniqueness(location *yaml.Node, items []NamedDefinition, errorMsg string) {
	itemsMap := map[string][]string{}
	for index := range items {
		item := items[index]
		itemId := item.Name.FlatCase()
		itemsGroup, found := itemsMap[itemId]
		if !found {
			itemsGroup = []string{}
			itemsMap[itemId] = itemsGroup
		}
		itemsMap[itemId] = append(itemsGroup, item.Name.Source)
	}
	for _, names := range itemsMap {
		if len(names) > 1 {
			validator.addError(location, fmt.Sprintf(`%s: %s`, errorMsg, strings.Join(names, ", ")))
		}
	}
}

func (validator *validator) EnumItemsUniqueness(location *yaml.Node, items []NamedEnumItem, errorMsg string) {
	itemsMap := map[string][]string{}
	for index := range items {
		item := items[index]
		itemId := item.Name.FlatCase()
		itemsGroup, found := itemsMap[itemId]
		if !found {
			itemsGroup = []string{}
			itemsMap[itemId] = itemsGroup
		}
		itemsMap[itemId] = append(itemsGroup, item.Name.Source)
	}
	for _, names := range itemsMap {
		if len(names) > 1 {
			validator.addError(location, fmt.Sprintf(`%s: %s`, errorMsg, strings.Join(names, ", ")))
		}
	}
}

func (validator *validator) NonEmpty(definition *NamedDefinition) {
	if definition.Type.Definition.IsEmpty() {
		validator.addError(definition.Location, "type empty can not be used in models")
	}
}

func (validator *validator) Model(model *NamedModel) {
	if model.IsObject() {
		validator.ItemsUniqueness(model.Location, model.Object.Fields, fmt.Sprintf(`object model %s fields names are too similiar to each other`, model.Name.Source))
		for index := range model.Object.Fields {
			field := model.Object.Fields[index]
			validator.NonEmpty(&field)
			validator.Definition(&field.Definition)
		}
	}
	if model.IsOneOf() {
		validator.ItemsUniqueness(model.Location, model.OneOf.Items, fmt.Sprintf(`oneOf model %s items names are too similiar to each other`, model.Name.Source))
		for index := range model.OneOf.Items {
			item := model.OneOf.Items[index]
			validator.NonEmpty(&item)
			validator.Definition(&item.Definition)
		}
	}
	if model.IsEnum() {
		validator.EnumItemsUniqueness(model.Location, model.Enum.Items, fmt.Sprintf(`enum model %s items names are too similiar to each other`, model.Name.Source))
	}
}

func (validator *validator) DefinitionDefault(definition *DefinitionDefault) {
	if definition != nil {
		if definition.Default != nil && !definition.Type.Definition.Info.Defaultable {
			validator.addError(definition.Location, fmt.Sprintf("type %s can not have default value", definition.Type.Definition.Name))
		}
		if definition.Default != nil {
			validator.DefaultValue(definition.Type.Definition, *definition.Default, definition.Location)
		}
	}
}

func (validator *validator) DefaultValue(typ TypeDef, value string, location *yaml.Node) {
	switch typ.Node {
	case ArrayType:
		if value != "[]" {
			validator.addError(location, fmt.Sprintf("default value for array type %s can be only empty list: [], found '%s'", typ.Name, value))
		}
	case MapType:
		if value != "{}" {
			validator.addError(location, fmt.Sprintf("default value for map type %s can be only empty map: {}, found '%s'", typ.Name, value))
		}
	case PlainType:
		switch typ.Plain {
		case TypeInt32,
			TypeInt64:
			err := Integer.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		case TypeDouble,
			TypeFloat,
			TypeDecimal:
			err := Float.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		case TypeBoolean:
			err := Boolean.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		case TypeUuid:
			err := UUID.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		case TypeDate:
			err := Date.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		case TypeDateTime:
			err := DateTime.Check(value)
			if err != nil {
				validator.addError(location, "default value "+err.Error())
			}
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				if !enumContainsItem(model.Enum, value) {
					validator.addError(location, fmt.Sprintf("default value %s is not defined in the enum %s", value, typ.Name))
				}
			}
		}
	}
}

func enumContainsItem(enum *Enum, what string) bool {
	for _, item := range enum.Items {
		if item.Name.Source == what {
			return true
		}
	}
	return false
}

func (validator *validator) Definition(definition *Definition) {
}
