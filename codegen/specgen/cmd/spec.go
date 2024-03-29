package cmd

import (
	"generator/console"
	"io/ioutil"
	"sort"
	"spec"
)

func readSpecFile(specFile string) *spec.Spec {
	console.PrintLnF("Reading spec file: %s", specFile)
	data, err := ioutil.ReadFile(specFile)
	FailIfErrorF(err, "Failed to read spec file: %s", specFile)

	console.PrintLn("Parsing spec")
	specification, messages, err := spec.ReadSpec(data)

	if messages != nil {
		sort.Sort(messages.Items)

		for _, message := range messages.Items {
			if message.Level != spec.LevelError {
				console.PrintLnF("%s %s", message.Level, message)
			} else {
				console.ProblemLnF("%s %s", message.Level, message)
			}
		}
	}

	FailIfErrorF(err, "Failed to parse spec: %s", specFile)
	return specification
}
