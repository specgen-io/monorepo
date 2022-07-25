package service

import "github.com/specgen-io/specgen/spec/v2"

func apiRouterName(api *spec.Api) string {
	return api.Name.CamelCase() + "Router"
}

func apiRouterNameVersioned(api *spec.Api) string {
	result := apiRouterName(api)
	version := api.Http.Version.Version
	if version.Source != "" {
		result = result + version.PascalCase()
	}
	return result
}

func apiServiceParamName(api *spec.Api) string {
	version := api.Http.Version
	name := api.Name.CamelCase() + "Service"
	if version.Version.Source != "" {
		name = name + version.Version.PascalCase()
	}
	return name
}
