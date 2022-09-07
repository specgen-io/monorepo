package spec

import (
	"fmt"
)

func OrderModels(models ModelsMap) []*NamedModel {
	walk := &walker{models, make(map[string]bool), nil}
	walk.models(models)
	return walk.orderedModels
}

type walker struct {
	modelsMap     ModelsMap
	visitedModels map[string]bool
	orderedModels []*NamedModel
}

func (walk *walker) models(models ModelsMap) {
	for _, model := range models {
		walk.model(model)
	}
}

func (walk *walker) model(model *NamedModel) {
	if _, visited := walk.visitedModels[model.Name.Source]; !visited {
		walk.visitedModels[model.Name.Source] = true
		if model.IsObject() {
			for index := range model.Object.Fields {
				field := model.Object.Fields[index]
				walk.typ(&field.Definition.Type.Definition)
			}
		}
		if model.IsOneOf() {
			for index := range model.OneOf.Items {
				item := model.OneOf.Items[index]
				walk.typ(&item.Definition.Type.Definition)
			}
		}
		walk.orderedModels = append(walk.orderedModels, model)
	}
}

func (walk *walker) typ(typ *TypeDef) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, found := walk.modelsMap[typ.Plain]; found {
				walk.model(model)
			}
		case NullableType:
		case ArrayType:
		case MapType:
			walk.typ(typ.Child)
		default:
			panic(fmt.Sprintf("unknown type node for type: %v", typ))
		}
	}
}
