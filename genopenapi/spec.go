package genopenapi

import "github.com/specgen-io/spec"

type UrlOperations struct {
	Url        string
	Operations []*spec.NamedOperation
}

func OperationsByUrl(specification *spec.Spec) []*UrlOperations {
	groups := make([]*UrlOperations, 0)
	groupsMap := make(map[string]*UrlOperations)
	for verIndex := range specification.Versions {
		version := &specification.Versions[verIndex]
		for apiIndex := range version.Http.Apis {
			api := &version.Http.Apis[apiIndex]
			for opIndex := range api.Operations {
				operation := &api.Operations[opIndex]
				url := operation.FullUrl()
				if _, contains := groupsMap[url]; !contains {
					urlOperations := UrlOperations{Url: url, Operations: make([]*spec.NamedOperation, 0)}
					groups = append(groups, &urlOperations)
					groupsMap[url] = &urlOperations
				}
				group, _ := groupsMap[url]
				group.Operations = append(group.Operations, operation)
			}
		}
	}
	return groups
}

func createEmptyResponse(name string, description string) *spec.NamedResponse {
	return &spec.NamedResponse{
		Name: spec.Name{name, nil},
		Definition: spec.Definition{
			Type:        spec.Type{Definition: spec.ParseType(spec.TypeEmpty), Location: nil},
			Description: &description,
			Location:    nil,
		},
	}
}

func addSpecialResponses(operation *spec.Operation) spec.Responses {
	responses := operation.Responses
	if operation.HasParams() || operation.Body != nil {
		if operation.GetResponse(spec.HttpStatusBadRequest) == nil {
			responses = append(responses, *createEmptyResponse(spec.HttpStatusBadRequest, "Service will return this if parameters are not provided or couldn't be parsed correctly"))
		}
	}
	if operation.GetResponse(spec.HttpStatusInternalServerError) == nil {
		responses = append(responses, *createEmptyResponse(spec.HttpStatusInternalServerError, "Service will return this if unexpected internal error happens"))
	}
	return responses
}
