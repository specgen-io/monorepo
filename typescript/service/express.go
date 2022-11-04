package service

import (
	"fmt"
	"strings"

	"generator"
	"spec"
	"typescript/modules"
	"typescript/validations"
	"typescript/validations/common"
	"typescript/writer"
)

type expressGenerator struct {
	validation validations.Validation
}

func (g *expressGenerator) SpecRouter(specification *spec.Spec, rootModule modules.Module, specRouterModule modules.Module) *generator.CodeFile {
	w := writer.New(specRouterModule)
	w.Line("import {Router} from 'express'")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			versionModule := rootModule.Submodule(version.Name.FlatCase())
			apiModule := versionModule.Submodule(serviceName(&api)) //TODO: This logic is repeated here, it also exists where api module is created
			routerModule := versionModule.Submodule("routing")      //TODO: This logic is repeated here, it also exists where router module is created
			w.Line("import {%s as %s} from '%s'", serviceInterfaceName(&api), serviceInterfaceNameVersioned(&api), apiModule.GetImport(specRouterModule))
			w.Line("import {%s as %s} from '%s'", apiRouterName(&api), apiRouterNameVersioned(&api), routerModule.GetImport(specRouterModule))
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
	w.Line("  const router = Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line("  router.use('%s', %s(services.%s))", expressVersionUrl(&version), apiRouterNameVersioned(&api), apiServiceParamName(&api))
		}
	}
	w.Line("  return router")
	w.Line("}")
	return w.ToCodeFile()
}

func expressVersionUrl(version *spec.Version) string {
	url := version.Http.GetUrl()
	if url == "" {
		return "/"
	}
	return url
}

func (g *expressGenerator) VersionRouting(version *spec.Version, targetModule modules.Module, modelsModule, validationModule, paramsModule, errorsModule, responsesModule modules.Module) *generator.CodeFile {
	w := writer.New(targetModule)

	w.Line(`import {Router, Request, Response} from 'express'`)
	w.Line(`import {zipHeaders} from '%s'`, paramsModule.GetImport(targetModule))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(targetModule))
	w.Line(`import * as models from '%s'`, modelsModule.GetImport(targetModule))
	w.Line(`import * as errors from '%s'`, errorsModule.GetImport(targetModule))
	w.Line(`import * as responses from '%s'`, responsesModule.GetImport(targetModule))

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
	return w.ToCodeFile()
}

func (g *expressGenerator) apiRouting(w generator.Writer, api *spec.Api) {
	w.EmptyLine()
	w.Line("export const %s = (service: %s) => {", apiRouterName(api), serviceInterfaceName(api))
	w.Line("  const router = Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.operationRouting(w.Indented(), &operation)
	}
	w.EmptyLine()
	w.Line("  return router")
	w.Line("}")
}

func getExpressUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), ":"+param.Name.Source, -1)
	}
	return url
}

func (g *expressGenerator) response(w generator.Writer, response *spec.Response, dataParam string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line("response.status(%s).send()", spec.HttpStatusCode(response.Name))
		w.Line("return")
	}
	if response.BodyIs(spec.BodyString) {
		w.Line("response.status(%s).type('text').send(%s)", spec.HttpStatusCode(response.Name), dataParam)
		w.Line("return")
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line("response.status(%s).type('json').send(JSON.stringify(t.encode(%s, %s)))", spec.HttpStatusCode(response.Name), g.validation.RuntimeType(&response.Type.Definition), dataParam)
		w.Line("return")
	}
}

func (g *expressGenerator) responses(w generator.Writer, responses spec.OperationResponses) {
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

func (g *expressGenerator) checkContentType(w generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`if (!responses.assertContentType(request, response, "text/plain")) {`)
		w.Line(`  return`)
		w.Line(`}`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`if (!responses.assertContentType(request, response, "application/json")) {`)
		w.Line(`  return`)
		w.Line(`}`)
	}
}

func (g *expressGenerator) operationRouting(w generator.Writer, operation *spec.NamedOperation) {
	w.Line("router.%s('%s', async (request: Request, response: Response) => {", strings.ToLower(operation.Endpoint.Method), getExpressUrl(operation.Endpoint))
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
	g.respondInternalServerError(w.Indented())
	w.Line("}")
	w.Unindent()
	w.Line("})")
}

func (g *expressGenerator) urlParamsParsing(w generator.Writer, operation *spec.NamedOperation) {
	if len(operation.Endpoint.UrlParams) > 0 {
		w.Line("const urlParamsDecode = t.decodeR(%s, request.params)", common.ParamsRuntimeTypeName(paramsTypeName(operation, "UrlParams")))
		w.Line("if (urlParamsDecode.error) {")
		g.respondNotFound(w.Indented(), "Failed to parse url parameters")
		w.Line("}")
		w.Line("const urlParams = urlParamsDecode.value")
	}
}

func (g *expressGenerator) headerParsing(w generator.Writer, operation *spec.NamedOperation) {
	if len(operation.HeaderParams) > 0 {
		w.Line("const headerParamsDecode = t.decodeR(%s, zipHeaders(request.rawHeaders))", common.ParamsRuntimeTypeName(paramsTypeName(operation, "HeaderParams")))
		w.Line("if (headerParamsDecode.error) {")
		g.respondBadRequest(w.Indented(), "HEADER", "headerParamsDecode.error", "Failed to parse header")
		w.Line("}")
		w.Line("const headerParams = headerParamsDecode.value")
	}
}

func (g *expressGenerator) queryParsing(w generator.Writer, operation *spec.NamedOperation) {
	if len(operation.QueryParams) > 0 {
		w.Line("const queryParamsDecode = t.decodeR(%s, request.query)", common.ParamsRuntimeTypeName(paramsTypeName(operation, "QueryParams")))
		w.Line("if (queryParamsDecode.error) {")
		g.respondBadRequest(w.Indented(), "QUERY", "queryParamsDecode.error", "Failed to parse query")
		w.Line("}")
		w.Line("const queryParams = queryParamsDecode.value")
	}
}

func (g *expressGenerator) bodyParsing(w generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`const body: string = request.body`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line("const bodyDecode = t.decodeR(%s, request.body)", g.validation.RuntimeType(&operation.Body.Type.Definition))
		w.Line("if (bodyDecode.error) {")
		g.respondBadRequest(w.Indented(), "BODY", "bodyDecode.error", "Failed to parse body")
		w.Line("}")
		w.Line("const body = bodyDecode.value")
	}
}

func (g *expressGenerator) respondBadRequest(w generator.Writer, location, errorsVar, message string) {
	w.Line(`responses.badRequest(response, { message: "%s", location: errors.ErrorLocation.%s, errors: %s })`, message, location, errorsVar)
	w.Line(`return`)
}

func (g *expressGenerator) respondNotFound(w generator.Writer, message string) {
	w.Line(`responses.notFound(response, { message: "%s" })`, message)
	w.Line(`return`)
}

func (g *expressGenerator) respondInternalServerError(w generator.Writer) {
	w.Line(`responses.internalServerError(response, { message: error instanceof Error ? error.message : "Unknown error" })`)
	w.Line(`return`)
}

func (g *expressGenerator) Responses(responsesModule, validationModule, errorsModule modules.Module) *generator.CodeFile {
	w := writer.New(responsesModule)

	w.Line(`import {Request, Response} from 'express'`)
	w.Line(`import * as t from '%s'`, validationModule.GetImport(responsesModule))
	w.Line(`import * as errors from '%s'`, errorsModule.GetImport(responsesModule))

	w.EmptyLine()
	code := `
export const internalServerError = (response: Response, error: errors.InternalServerError) => {
  const body = t.encode(errors.TInternalServerError, error)
  response.status(500).type("json").send(JSON.stringify(body))
}

export const notFound = (response: Response, error: errors.NotFoundError) => {
  const body = t.encode(errors.TNotFoundError, error)
  response.status(404).type("json").send(JSON.stringify(body))
}

export const badRequest = (response: Response, error: errors.BadRequestError) => {
  const body = t.encode(errors.TBadRequestError, error)
  response.status(400).type("json").send(JSON.stringify(body))
}

export const assertContentType = (request: Request, response: Response, contentType: string): boolean => {
  if (!request.is(contentType)) {
    const error = {
      message: "Failed to parse header", 
      location: errors.ErrorLocation.HEADER,
      errors: [{path: "Content-Type", code: "missing", message: 'Expected Content-Type header: ${contentType}'}]
    }
    badRequest(response, error)
    return false
  }
  return true
}
`
	code = strings.Replace(code, "'", "`", -1)
	w.Lines(code)

	return w.ToCodeFile()
}
