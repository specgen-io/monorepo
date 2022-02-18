package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/spec"
	"io/ioutil"
	"sort"
)

func readSpecFile(specFile string) *spec.Spec {
	console.PrintLnF("Reading spec file: %s", specFile)
	data, err := ioutil.ReadFile(specFile)
	fail.IfErrorF(err, "Failed to read spec file: %s", specFile)

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

	fail.IfErrorF(err, "Failed to parse spec: %s", specFile)
	return specification
}
