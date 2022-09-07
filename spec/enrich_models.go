package spec

import "fmt"

func enrichModels(models Models, messages *Messages) []*NamedModel {
	enricher := &modelsEnricher{buildModelsMap(models), make(map[string]bool), messages, nil}
	for index := range models {
		enricher.model(&models[index])
	}
	return enricher.orderedModels
}

type modelsEnricher struct {
	models        ModelsMap
	visitedModels map[string]bool
	messages      *Messages
	orderedModels []*NamedModel
}

func (enricher *modelsEnricher) model(model *NamedModel) {
	if _, visited := enricher.visitedModels[model.Name.Source]; !visited {
		enricher.visitedModels[model.Name.Source] = true
		if model.IsObject() {
			for index := range model.Object.Fields {
				field := &model.Object.Fields[index]
				enricher.typ(&field.Definition.Type)
			}
		}
		if model.IsOneOf() {
			for index := range model.OneOf.Items {
				item := &model.OneOf.Items[index]
				enricher.typ(&item.Definition.Type)
			}
		}
		enricher.orderedModels = append(enricher.orderedModels, model)
	}
}

func (enricher *modelsEnricher) typ(typ *Type) {
	enricher.typeDef(typ, &typ.Definition)
}

func (enricher *modelsEnricher) typeDef(starter *Type, typ *TypeDef) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, found := enricher.models[typ.Plain]; found {
				typ.Info = ModelTypeInfo(model)
				enricher.model(model)
			} else {
				if info, ok := Types[typ.Plain]; ok {
					typ.Info = &info
				} else {
					e := Error("unknown type: %s", typ.Plain).At(locationFromNode(starter.Location))
					enricher.messages.Add(e)
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
