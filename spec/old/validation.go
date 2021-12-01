package old

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"strconv"
)

const (
	LevelError   string = "error"
	LevelWarning string = "warning"
)

type Message struct {
	Level    string
	Message  string
	Location *yaml.Node
}

type Messages []Message

func (self Message) String() string {
	lineStr := "unknown"
	if self.Location != nil {
		lineStr = strconv.Itoa(self.Location.Line)
	}
	return fmt.Sprintf("line %s: %s", lineStr, self.Message)
}

type validator struct {
	Errors   Messages
	Warnings Messages
}

func validate(old *Spec) (Messages, Messages) {
	validator := &validator{}
	validator.Spec(old)
	return validator.Warnings, validator.Errors
}

func (validator *validator) AddError(node *yaml.Node, message string) {
	theMessage := Message{Level: LevelError, Message: message, Location: node}
	validator.Errors = append(validator.Errors, theMessage)
}

func (validator *validator) AddWarning(node *yaml.Node, message string) {
	theMessage := Message{Level: LevelWarning, Message: message, Location: node}
	validator.Warnings = append(validator.Warnings, theMessage)
}

func (validator *validator) Spec(old *Spec) {
	for versionIndex := range old.Versions {
		models := old.Versions[versionIndex].Models
		for modelIndex := range models {
			validator.Model(&models[modelIndex])
		}
		apis := old.Versions[versionIndex].Http.Apis
		for apiIndex := range apis {
			for operation := range apis[apiIndex].Operations {
				validator.Operation(&apis[apiIndex].Operations[operation])
			}
		}
	}
}

func (validator *validator) ParamsNames(paramsMap map[string]NamedParam, params []NamedParam) {
	for _, p := range params {
		if other, ok := paramsMap[p.Name.SnakeCase()]; ok {
			message := fmt.Sprintf("parameter name '%s' conflicts with the other parameter name '%s'", p.Name.Source, other.Name.Source)
			validator.AddError(p.Name.Location, message)
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
	validator.Params(operation.HeaderParams, false)

	if operation.Body != nil && !operation.Body.Type.Definition.IsEmpty() {
		if operation.Body.Type.Definition.Info.Structure != StructureObject && operation.Body.Type.Definition.Info.Structure != StructureArray {
			message := fmt.Sprintf("body should be of a type with structure of an object or array, found %s", operation.Body.Type.Definition.Name)
			validator.AddError(operation.Body.Location, message)
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
			validator.AddWarning(response.Type.Location, message)
		}
	}
	if !response.Type.Definition.IsEmpty() && response.Type.Definition.Info.Structure != StructureObject && response.Type.Definition.Info.Structure != StructureArray {
		message := fmt.Sprintf("response %s should be either empty or some type with structure of an object or array, found %s", response.Name.Source, response.Type.Definition.Name)
		validator.AddError(response.Type.Location, message)
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
				validator.AddError(paramType.Location, fmt.Sprintf("parameter %s should be of scalar type or array of scalar type, found %s", paramName.Source, paramType.Definition.Name))
			}
		} else {
			if !scalar {
				validator.AddError(paramType.Location, fmt.Sprintf("parameter %s should be of scalar type, found %s", paramName.Source, paramType.Definition.Name))
			}
		}
		validator.DefinitionDefault(&params[index].DefinitionDefault)
	}
}

func (validator *validator) Model(model *NamedModel) {
	if model.IsObject() {
		for index := range model.Object.Fields {
			validator.Definition(&model.Object.Fields[index].Definition)
		}
	}
}

func (validator *validator) DefinitionDefault(definition *DefinitionDefault) {
	if definition != nil {
		if definition.Default != nil && !definition.Type.Definition.Info.Defaultable {
			validator.AddError(definition.Location, fmt.Sprintf("type %s can not have default value", definition.Type.Definition.Name))
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
			validator.AddError(location, fmt.Sprintf("default value for array type %s can be only empty list: [], found '%s'", typ.Name, value))
		}
	case MapType:
		if value != "{}" {
			validator.AddError(location, fmt.Sprintf("default value for map type %s can be only empty map: {}, found '%s'", typ.Name, value))
		}
	case PlainType:
		switch typ.Plain {
		case TypeInt32,
			TypeInt64:
			err := Integer.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDouble,
			TypeFloat,
			TypeDecimal:
			err := Float.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeBoolean:
			err := Boolean.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeUuid:
			err := UUID.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDate:
			err := Date.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		case TypeDateTime:
			err := DateTime.Check(value)
			if err != nil {
				validator.AddError(location, "default value "+err.Error())
			}
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				if !enumContainsItem(model.Enum, value) {
					validator.AddError(location, fmt.Sprintf("default value %s is not defined in the enum %s", value, typ.Name))
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
	if definition != nil {
	}
}