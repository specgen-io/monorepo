// Generated by ego.
// DO NOT EDIT

//line gents/iotsmodels.ego:1
package gents

import
//line gents/iotsmodels.ego:2

//line gents/iotsmodels.ego:3
"fmt"
import "html"
import "io"
import "context"

import (
//line gents/iotsmodels.ego:4

//line gents/iotsmodels.ego:5
	spec "github.com/specgen-io/spec.v1"

//line gents/iotsmodels.ego:6
)

//line gents/iotsmodels.ego:7

//line gents/iotsmodels.ego:8

//line gents/iotsmodels.ego:9
func generateIoTsObjectModel(model spec.NamedModel, w io.Writer) {
//line gents/iotsmodels.ego:10
	hasRequiredFields, hasOptionalFields := kindOfFields(model)
//line gents/iotsmodels.ego:12
	_, _ = io.WriteString(w, "\n")
//line gents/iotsmodels.ego:12
	if hasRequiredFields && hasOptionalFields {
//line gents/iotsmodels.ego:13
		_, _ = io.WriteString(w, "export const T")
//line gents/iotsmodels.ego:13
		_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:13
		_, _ = io.WriteString(w, " = t.intersection([\n  t.interface({\n")
//line gents/iotsmodels.ego:15
		for _, field := range model.Object.Fields {
//line gents/iotsmodels.ego:16
			if !field.Type.Definition.IsNullable() {
//line gents/iotsmodels.ego:17
				_, _ = io.WriteString(w, "    ")
//line gents/iotsmodels.ego:17
				_, _ = fmt.Fprint(w, field.Name.Source)
//line gents/iotsmodels.ego:17
				_, _ = io.WriteString(w, ": ")
//line gents/iotsmodels.ego:17
				_, _ = fmt.Fprint(w, IoTsType(&field.Type.Definition))
//line gents/iotsmodels.ego:17
				_, _ = io.WriteString(w, ",\n")
//line gents/iotsmodels.ego:18
			}
//line gents/iotsmodels.ego:19
		}
//line gents/iotsmodels.ego:20
		_, _ = io.WriteString(w, "  }),\n  t.partial({\n")
//line gents/iotsmodels.ego:22
		for _, field := range model.Object.Fields {
//line gents/iotsmodels.ego:23
			if field.Type.Definition.IsNullable() {
//line gents/iotsmodels.ego:24
				_, _ = io.WriteString(w, "    ")
//line gents/iotsmodels.ego:24
				_, _ = fmt.Fprint(w, field.Name.Source)
//line gents/iotsmodels.ego:24
				_, _ = io.WriteString(w, ": ")
//line gents/iotsmodels.ego:24
				_, _ = fmt.Fprint(w, IoTsType(&field.Type.Definition))
//line gents/iotsmodels.ego:24
				_, _ = io.WriteString(w, ",\n")
//line gents/iotsmodels.ego:25
			}
//line gents/iotsmodels.ego:26
		}
//line gents/iotsmodels.ego:27
		_, _ = io.WriteString(w, "  })\n])\n")
//line gents/iotsmodels.ego:29
	} else {
//line gents/iotsmodels.ego:30
		var ioTsType = "t.interface"
//line gents/iotsmodels.ego:31
		if hasOptionalFields {
//line gents/iotsmodels.ego:32
			ioTsType = "t.partial"
//line gents/iotsmodels.ego:33
		}
//line gents/iotsmodels.ego:34
		_, _ = io.WriteString(w, "export const T")
//line gents/iotsmodels.ego:34
		_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:34
		_, _ = io.WriteString(w, " = ")
//line gents/iotsmodels.ego:34
		_, _ = fmt.Fprint(w, ioTsType)
//line gents/iotsmodels.ego:34
		_, _ = io.WriteString(w, "({\n")
//line gents/iotsmodels.ego:35
		for _, field := range model.Object.Fields {
//line gents/iotsmodels.ego:36
			_, _ = io.WriteString(w, "  ")
//line gents/iotsmodels.ego:36
			_, _ = fmt.Fprint(w, field.Name.Source)
//line gents/iotsmodels.ego:36
			_, _ = io.WriteString(w, ": ")
//line gents/iotsmodels.ego:36
			_, _ = fmt.Fprint(w, IoTsType(&field.Type.Definition))
//line gents/iotsmodels.ego:36
			_, _ = io.WriteString(w, ",\n")
//line gents/iotsmodels.ego:37
		}
//line gents/iotsmodels.ego:38
		_, _ = io.WriteString(w, "})\n")
//line gents/iotsmodels.ego:39
	}
//line gents/iotsmodels.ego:41
	_, _ = io.WriteString(w, "\nexport type ")
//line gents/iotsmodels.ego:41
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:41
	_, _ = io.WriteString(w, " = t.TypeOf<typeof T")
//line gents/iotsmodels.ego:41
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:41
	_, _ = io.WriteString(w, ">\n")
//line gents/iotsmodels.ego:42
}

//line gents/iotsmodels.ego:43

//line gents/iotsmodels.ego:44

//line gents/iotsmodels.ego:45
func generateIoTsEnumModel(model spec.NamedModel, w io.Writer) {
//line gents/iotsmodels.ego:47
	_, _ = io.WriteString(w, "\nexport enum ")
//line gents/iotsmodels.ego:47
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:47
	_, _ = io.WriteString(w, " {\n")
//line gents/iotsmodels.ego:48
	for _, item := range model.Enum.Items {
//line gents/iotsmodels.ego:49
		_, _ = io.WriteString(w, "  ")
//line gents/iotsmodels.ego:49
		_, _ = fmt.Fprint(w, item.Name.UpperCase())
//line gents/iotsmodels.ego:49
		_, _ = io.WriteString(w, " = \"")
//line gents/iotsmodels.ego:49
		_, _ = fmt.Fprint(w, item.Value)
//line gents/iotsmodels.ego:49
		_, _ = io.WriteString(w, "\",\n")
//line gents/iotsmodels.ego:50
	}
//line gents/iotsmodels.ego:51
	_, _ = io.WriteString(w, "}\n\nexport const T")
//line gents/iotsmodels.ego:53
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:53
	_, _ = io.WriteString(w, " = t.enum(")
//line gents/iotsmodels.ego:53
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:53
	_, _ = io.WriteString(w, ")\n")
//line gents/iotsmodels.ego:54
}

//line gents/iotsmodels.ego:55

//line gents/iotsmodels.ego:56

//line gents/iotsmodels.ego:57
func generateIoTsUnionModel(model spec.NamedModel, w io.Writer) {
//line gents/iotsmodels.ego:59
	_, _ = io.WriteString(w, "\nexport const T")
//line gents/iotsmodels.ego:59
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:59
	_, _ = io.WriteString(w, " = t.union([\n")
//line gents/iotsmodels.ego:60
	for _, item := range model.OneOf.Items {
//line gents/iotsmodels.ego:61
		_, _ = io.WriteString(w, "  t.interface({")
//line gents/iotsmodels.ego:61
		_, _ = fmt.Fprint(w, item.Name.Source)
//line gents/iotsmodels.ego:61
		_, _ = io.WriteString(w, ": ")
//line gents/iotsmodels.ego:61
		_, _ = fmt.Fprint(w, IoTsType(&item.Type.Definition))
//line gents/iotsmodels.ego:61
		_, _ = io.WriteString(w, "}),\n")
//line gents/iotsmodels.ego:62
	}
//line gents/iotsmodels.ego:63
	_, _ = io.WriteString(w, "])\n\nexport type ")
//line gents/iotsmodels.ego:65
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:65
	_, _ = io.WriteString(w, " = t.TypeOf<typeof T")
//line gents/iotsmodels.ego:65
	_, _ = fmt.Fprint(w, model.Name.PascalCase())
//line gents/iotsmodels.ego:65
	_, _ = io.WriteString(w, ">\n")
//line gents/iotsmodels.ego:66
}

//line gents/iotsmodels.ego:67

//line gents/iotsmodels.ego:68

//line gents/iotsmodels.ego:69
func generateIoTsModels(spec *spec.Spec, w io.Writer) {
//line gents/iotsmodels.ego:70
	_, _ = io.WriteString(w, "/* eslint-disable @typescript-eslint/camelcase */\n/* eslint-disable @typescript-eslint/no-magic-numbers */\nimport * as t from './io-ts'\n")
//line gents/iotsmodels.ego:73
	for _, model := range spec.ResolvedModels {
//line gents/iotsmodels.ego:74
		if model.IsObject() {
//line gents/iotsmodels.ego:75
			generateIoTsObjectModel(model, w)
//line gents/iotsmodels.ego:76
		} else if model.IsEnum() {
//line gents/iotsmodels.ego:77
			generateIoTsEnumModel(model, w)
//line gents/iotsmodels.ego:78
		} else if model.IsOneOf() {
//line gents/iotsmodels.ego:79
			generateIoTsUnionModel(model, w)
//line gents/iotsmodels.ego:80
		}
//line gents/iotsmodels.ego:81
	}
//line gents/iotsmodels.ego:82
}

var _ fmt.Stringer
var _ io.Reader
var _ context.Context
var _ = html.EscapeString
