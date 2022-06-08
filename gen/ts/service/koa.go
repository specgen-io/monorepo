package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/types"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

type koaGenerator struct {
	validation validations.Validation
}

func (g *koaGenerator) SpecRouter(specification *spec.Spec, rootModule modules.Module, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	w.Line("import Router from '@koa/router'")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			versionModule := rootModule.Submodule(version.Version.FlatCase())
			apiModule := versionModule.Submodule(serviceName(&api)) //TODO: This logic is repeated here, it also exists where api module is created
			routerModule := versionModule.Submodule("routing")      //TODO: This logic is repeated here, it also exists where router module is created
			w.Line("import {%s as %s} from '%s'", serviceInterfaceName(&api), serviceInterfaceNameVersioned(&api), apiModule.GetImport(module))
			w.Line("import {%s as %s} from '%s'", apiRouterName(&api), apiRouterNameVersioned(&api), routerModule.GetImport(module))
		}
	}

	servicesDefinitions := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			servicesDefinitions = append(servicesDefinitions, fmt.Sprintf("%s: %s", apiServiceParamName(&api), serviceInterfaceNameVersioned(&api)))
		}
	}

	w.EmptyLine()
	w.Line(`export interface Services {`)
	for _, serviceDefinition := range servicesDefinitions {
		w.Line(`  %s`, serviceDefinition)
	}
	w.Line(`}`)

	w.EmptyLine()
	w.Line("export const specRouter = (services: Services) => {")
	w.Line("  const router = new Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiRouterNameConst := "the" + casee.ToPascalCase(apiRouterNameVersioned(&api))
			w.Line("  const %s = %s(services.%s)", apiRouterNameConst, apiRouterNameVersioned(&api), apiServiceParamName(&api))
			if version.Http.GetUrl() == "" {
				w.Line("  router.use(%s.routes(), %s.allowedMethods())", apiRouterNameConst, apiRouterNameConst)
			} else {
				w.Line("  router.use('%s', %s.routes(), %s.allowedMethods())", version.Http.GetUrl(), apiRouterNameConst, apiRouterNameConst)
			}
		}
	}
	w.Line("  return router")
	w.Line("}")

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func (g *koaGenerator) VersionRouting(version *spec.Version, validationModule, paramsModule, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()

	w.Line(`import Router from '@koa/router'`)
	w.Line(`import { ExtendableContext } from 'koa'`)
	w.Line(`import {zipHeaders} from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as models from './models'`)
	for _, api := range version.Http.Apis {
		w.Line("import {%s} from './%s'", serviceInterfaceName(&api), serviceName(&api))
	}

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "HeaderParams"), operation.HeaderParams)
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "UrlParams"), operation.Endpoint.UrlParams)
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "QueryParams"), operation.QueryParams)
		}

		g.apiRouting(w, &api)
	}

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func (g *koaGenerator) apiRouting(w *sources.Writer, api *spec.Api) {
	w.EmptyLine()
	w.Line("export const %s = (service: %s) => {", apiRouterName(api), serviceInterfaceName(api))
	g.generateErrors(w.Indented())
	w.EmptyLine()
	w.Line("  const router = new Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.operationRouting(w.Indented(), &operation)
	}
	w.EmptyLine()
	w.Line("  return router")
	w.Line("}")
}

func getKoaUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), ":"+param.Name.Source, -1)
	}
	return url
}

func (g *koaGenerator) response(w *sources.Writer, response *spec.Response, dataParam string) {
	w.Line("ctx.status = %s", spec.HttpStatusCode(response.Name))
	if response.BodyIs(spec.BodyEmpty) {
		w.Line("return")
	}
	if response.BodyIs(spec.BodyString) {
		w.Line("ctx.body = %s", dataParam)
		w.Line("return")
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line("ctx.body = t.encode(%s, %s)", g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &response.Type.Definition), dataParam)
		w.Line("return")
	}
}

func (g *koaGenerator) responses(w *sources.Writer, responses spec.Responses) {
	if len(responses) == 1 {
		g.response(w, &responses[0].Response, "result")
	} else {
		w.Line("switch (result.status) {")
		for _, response := range responses {
			w.Line("  case '%s':", response.Name.SnakeCase())
			g.response(w.IndentedWith(2), &response.Response, "result.data")
		}
		w.Line("}")
	}
}

func (g *koaGenerator) checkContentType(w *sources.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`if (!assertContentType(ctx, "text/plain")) {`)
		w.Line(`  return`)
		w.Line(`}`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`if (!assertContentType(ctx, "application/json")) {`)
		w.Line(`  return`)
		w.Line(`}`)
	}
}

func (g *koaGenerator) operationRouting(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line("router.%s('%s', async (ctx) => {", strings.ToLower(operation.Endpoint.Method), getKoaUrl(operation.Endpoint))
	w.Indent()
	w.Line("try {")
	w.Indent()
	g.urlParamsParsing(w, operation)
	g.checkContentType(w, operation)
	g.headerParsing(w, operation)
	g.queryParsing(w, operation)
	g.bodyParsing(w, operation)
	w.Line(serviceCall(operation, getApiCallParamsObject(operation)))
	g.responses(w, operation.Responses)
	w.Unindent()
	w.Line("} catch (error) {")
	g.respondInternalServerError(w.Indented(), operation)
	w.Line("}")
	w.Unindent()
	w.Line("})")
}

func (g *koaGenerator) urlParamsParsing(w *sources.Writer, operation *spec.NamedOperation) {
	if len(operation.Endpoint.UrlParams) > 0 {
		w.Line("const urlParamsDecode = t.decodeR(%s, ctx.params)", urlParamsRuntimeType(operation))
		w.Line("if (urlParamsDecode.error) {")
		g.respondNotFound(w.Indented(), operation, "Failed to parse url parameters")
		w.Line("}")
		w.Line("const urlParams = urlParamsDecode.value")
	}
}

func (g *koaGenerator) headerParsing(w *sources.Writer, operation *spec.NamedOperation) {
	if len(operation.HeaderParams) > 0 {
		w.Line("const headerParamsDecode = t.decodeR(%s, zipHeaders(ctx.req.rawHeaders))", headersRuntimeType(operation))
		w.Line("if (headerParamsDecode.error) {")
		g.respondBadRequest(w.Indented(), operation, "HEADER", "headerParamsDecode.error", "Failed to parse header parameters")
		w.Line("}")
		w.Line("const headerParams = headerParamsDecode.value")
	}
}

func (g *koaGenerator) queryParsing(w *sources.Writer, operation *spec.NamedOperation) {
	if len(operation.QueryParams) > 0 {
		w.Line("const queryParamsDecode = t.decodeR(%s, ctx.request.query)", queryRuntimeType(operation))
		w.Line("if (queryParamsDecode.error) {")
		g.respondBadRequest(w.Indented(), operation, "QUERY", "queryParamsDecode.error", "Failed to parse query parameters")
		w.Line("}")
		w.Line("const queryParams = queryParamsDecode.value")
	}
}

func (g *koaGenerator) bodyParsing(w *sources.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`const body: string = ctx.request.rawBody`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line("const bodyDecode = t.decodeR(%s, ctx.request.body)", g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &operation.Body.Type.Definition))
		w.Line("if (bodyDecode.error) {")
		g.respondBadRequest(w.Indented(), operation, "BODY", "bodyDecode.error", "Failed to parse body JSON")
		w.Line("}")
		w.Line("const body = bodyDecode.value")
	}
}

func (g *koaGenerator) respondBadRequest(w *sources.Writer, operation *spec.NamedOperation, location, errorsVar, message string) {
	w.Line(`respondBadRequest(ctx, { message: "%s", location: models.ErrorLocation.%s, errors: %s })`, message, location, errorsVar)
	w.Line(`return`)
}

func (g *koaGenerator) respondNotFound(w *sources.Writer, operation *spec.NamedOperation, message string) {
	w.Line(`respondNotFound(ctx, { message: "%s" })`, message)
	w.Line(`return`)
}

func (g *koaGenerator) respondInternalServerError(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line(`respondInternalServerError(ctx, { message: error instanceof Error ? error.message : "Unknown error" })`)
	w.Line(`return`)
}

func (g *koaGenerator) generateErrors(w *sources.Writer) {
	code := `
const respondInternalServerError = (ctx: ExtendableContext, error: models.InternalServerError) => {
  const body = t.encode(models.TInternalServerError, error)
  ctx.status = 500
  ctx.body = body
}

const respondNotFound = (ctx: ExtendableContext, error: models.NotFoundError) => {
  const body = t.encode(models.TNotFoundError, error)
  ctx.status = 404
  ctx.body = body
}

const respondBadRequest = (ctx: ExtendableContext, error: models.BadRequestError) => {
  const body = t.encode(models.TBadRequestError, error)
  ctx.status = 400
  ctx.body = body
}

const assertContentType = (ctx: ExtendableContext, contentType: string): boolean => {
  if (ctx.request.type != contentType) {
    const message = 'Expected Content-Type header: ${contentType}'
    const error = {
      message, 
      location: models.ErrorLocation.HEADER,
      errors: [{path: "Content-Type", code: "wrong_value", message}]
    }
    respondBadRequest(ctx, error)
    return false
  }
  return true
}
`
	w.Lines(strings.Replace(code, "'", "`", -1))
}
