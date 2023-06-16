package spec

import "fmt"

type SpecWalker struct {
	onSpecification     func(specification *Spec)
	onHttpErrors        func(httpErrors *HttpErrors)
	onResponse          func(response *Response)
	onVersion           func(version *Version)
	onApi               func(api *Api)
	onOperation         func(operation *NamedOperation)
	onOperationResponse func(response *OperationResponse)
	onParam             func(param *NamedParam)
	onModel             func(model *NamedModel)
	onType              func(typ *Type)
	onTypeDef           func(typ *TypeDef)
}

func NewWalker() *SpecWalker {
	return &SpecWalker{}
}

func (w *SpecWalker) OnSpecification(callback func(specification *Spec)) *SpecWalker {
	w.onSpecification = callback
	return w
}

func (w *SpecWalker) OnHttpErrors(callback func(httpErrors *HttpErrors)) *SpecWalker {
	w.onHttpErrors = callback
	return w
}

func (w *SpecWalker) OnResponse(callback func(response *Response)) *SpecWalker {
	w.onResponse = callback
	return w
}

func (w *SpecWalker) OnVersion(callback func(version *Version)) *SpecWalker {
	w.onVersion = callback
	return w
}

func (w *SpecWalker) OnApi(callback func(api *Api)) *SpecWalker {
	w.onApi = callback
	return w
}

func (w *SpecWalker) OnOperation(callback func(operation *NamedOperation)) *SpecWalker {
	w.onOperation = callback
	return w
}

func (w *SpecWalker) OnOperationResponse(callback func(response *OperationResponse)) *SpecWalker {
	w.onOperationResponse = callback
	return w
}

func (w *SpecWalker) OnParam(callback func(param *NamedParam)) *SpecWalker {
	w.onParam = callback
	return w
}

func (w *SpecWalker) OnModel(callback func(model *NamedModel)) *SpecWalker {
	w.onModel = callback
	return w
}

func (w *SpecWalker) OnType(callback func(typ *Type)) *SpecWalker {
	w.onType = callback
	return w
}

func (w *SpecWalker) OnTypeDef(callback func(typ *TypeDef)) *SpecWalker {
	w.onTypeDef = callback
	return w
}

func (w *SpecWalker) Specification(specification *Spec) {
	if w.onSpecification != nil {
		w.onSpecification(specification)
	}
	if specification.HttpErrors != nil {
		w.HttpErrors(specification.HttpErrors)
	}
	for index := range specification.Versions {
		w.Version(&specification.Versions[index])
	}
}

func (w *SpecWalker) HttpErrors(httpErrors *HttpErrors) {
	if w.onHttpErrors != nil {
		w.onHttpErrors(httpErrors)
	}
	for index := range httpErrors.Models {
		w.Model(&httpErrors.Models[index])
	}
	for index := range httpErrors.Responses {
		w.Response(&httpErrors.Responses[index].Response)
	}
}

func (w *SpecWalker) Models(models []*NamedModel) {
	for index := range models {
		w.Model(models[index])
	}
}

func (w *SpecWalker) Response(response *Response) {
	if w.onResponse != nil {
		w.onResponse(response)
	}
	w.Type(&response.Definition.Type)
}

func (w *SpecWalker) Version(version *Version) {
	if w.onVersion != nil {
		w.onVersion(version)
	}
	for index := range version.Models {
		w.Model(&version.Models[index])
	}
	for index := range version.Http.Apis {
		w.Api(&version.Http.Apis[index])
	}
}

func (w *SpecWalker) Api(api *Api) {
	if w.onApi != nil {
		w.onApi(api)
	}
	for index := range api.Operations {
		operation := &api.Operations[index]
		w.Operation(operation)
	}
}

func (w *SpecWalker) Operation(operation *NamedOperation) {
	if w.onOperation != nil {
		w.onOperation(operation)
	}
	w.params(operation.Endpoint.UrlParams)
	w.params(operation.QueryParams)
	w.params(operation.HeaderParams)

	if operation.Body != nil {
		w.Type(operation.Body.Type)
	}

	for index := range operation.Responses {
		w.OperationResponse(&operation.Responses[index])
	}
}

func (w *SpecWalker) OperationResponse(response *OperationResponse) {
	if w.onOperationResponse != nil {
		w.onOperationResponse(response)
	}
	w.Type(&response.Definition.Type)
}

func (w *SpecWalker) params(params []NamedParam) {
	for index := range params {
		w.Param(&params[index])
	}
}

func (w *SpecWalker) Param(param *NamedParam) {
	if w.onParam != nil {
		w.onParam(param)
	}
	w.Type(&param.DefinitionDefault.Type)
}

func (w *SpecWalker) Model(model *NamedModel) {
	if w.onModel != nil {
		w.onModel(model)
	}
	if model.IsObject() {
		for index := range model.Object.Fields {
			field := &model.Object.Fields[index]
			w.Type(&field.Definition.Type)
		}
	}
	if model.IsOneOf() {
		for index := range model.OneOf.Items {
			item := &model.OneOf.Items[index]
			w.Type(&item.Definition.Type)
		}
	}
}

func (w *SpecWalker) Type(typ *Type) {
	if w.onType != nil {
		w.onType(typ)
	}
	w.TypeDef(&typ.Definition)
}

func (w *SpecWalker) TypeDef(typ *TypeDef) {
	if typ != nil {
		if w.onTypeDef != nil {
			w.onTypeDef(typ)
		}
		switch typ.Node {
		case PlainType:
			return
		case NullableType:
		case ArrayType:
		case MapType:
			w.TypeDef(typ.Child)
		default:
			panic(fmt.Sprintf("unknown kind of type: %v", typ))
		}
	}
}
