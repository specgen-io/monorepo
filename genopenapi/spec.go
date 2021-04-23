package genopenapi

import spec "github.com/specgen-io/spec.v2"

type Version string

type Operation struct {
	Version   spec.Name
	Api       spec.Api
	Operation spec.NamedOperation
}

type Group struct {
	Url        string
	Operations []Operation
}

func Groups(versions []spec.VersionedApis) []*Group {
	groups := make([]*Group, 0)
	groupsMap := make(map[string]*Group)
	for _, version := range versions {
		for _, api := range version.Apis {
			for _, operation := range api.Operations {
				url := version.GetUrl() + operation.Endpoint.Url
				if _, contains := groupsMap[url]; !contains {
					group := Group{Url: url, Operations: make([]Operation, 0)}
					groups = append(groups, &group)
					groupsMap[url] = &group
				}
				group, _ := groupsMap[url]
				group.Operations = append(group.Operations, Operation{version.Version, api, operation})
			}
		}
	}

	return groups
}
