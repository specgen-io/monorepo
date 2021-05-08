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

func Groups(specification *spec.Spec) []*Group {
	groups := make([]*Group, 0)
	groupsMap := make(map[string]*Group)
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			for _, operation := range api.Operations {
				url := operation.FullUrl()
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
