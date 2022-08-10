package service

import "spec"

func apiRouterName(api *spec.Api) string {
	return api.Name.CamelCase() + "Router"
}

func apiRouterNameVersioned(api *spec.Api) string {
	result := apiRouterName(api)
	version := api.Http.Version.Name
	if version.Source != "" {
		result = result + version.PascalCase()
	}
	return result
}

func apiServiceParamName(api *spec.Api) string {
	version := api.Http.Version
	name := api.Name.CamelCase() + "Service"
	if version.Name.Source != "" {
		name = name + version.Name.PascalCase()
	}
	return name
}
