package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/spec"
	"io/ioutil"
	"sort"
)

func printSpecParseResult(messages spec.Messages) {
	if messages == nil {
		return
	}
	sort.Sort(messages)

	for _, message := range messages {
		if message.Level != spec.LevelError {
			console.PrintLnF("%s %s", message.Level, message)
		} else {
			console.ProblemLnF("%s %s", message.Level, message)
		}
	}
}

func readSpecFile(specFile string) *spec.Spec {
	console.PrintLnF("Reading spec file: %s", specFile)
	data, err := ioutil.ReadFile(specFile)
	fail.IfErrorF(err, "Failed to read spec file: %s", specFile)

	console.PrintLn("Parsing spec")
	spec, messages, err := spec.ReadSpec(data)

	printSpecParseResult(messages)
	fail.IfErrorF(err, "Failed to parse spec: %s", specFile)
	return spec
}
