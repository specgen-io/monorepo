package openapi

import "github.com/specgen-io/specgen/v2/spec"

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
