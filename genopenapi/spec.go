package genopenapi

import "github.com/specgen-io/spec"

type Operation struct {
	Api       spec.Api
	Operation spec.NamedOperation
}

type Group struct {
	Url        string
	Operations []Operation
}

func Groups(apis spec.Apis) []*Group {
	groups := make([]*Group, 0)
	groupsMap := make(map[string]*Group)

	for _, api := range apis {
		for _, operation := range api.Operations {
			url := operation.Endpoint.Url
			if _, contains := groupsMap[url]; !contains {
				group := Group{Url: url, Operations: make([]Operation, 0)}
				groups = append(groups, &group)
				groupsMap[url] = &group
			}
			group, _ := groupsMap[url]
			group.Operations = append(group.Operations, Operation{api, operation})
		}
	}

	return groups
}
