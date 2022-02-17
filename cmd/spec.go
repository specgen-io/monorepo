package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/spec"
	"io/ioutil"
	"sort"
)

type messages spec.Messages

func (ms messages) Len() int {
	return len(ms)
}
func (ms messages) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
func (ms messages) Less(i, j int) bool {
	if ms[i].Location == nil {
		return true
	}
	if ms[j].Location == nil {
		return false

	}
	return ms[i].Location.Line < ms[j].Location.Line ||
		(ms[i].Location.Line == ms[j].Location.Line && ms[i].Location.Column < ms[j].Location.Column)
}

func printSpecParseResult(result *spec.SpecParseResult) {
	allMessages := append(result.Errors, result.Warnings...)
	sort.Sort(messages(allMessages))

	for _, message := range allMessages {
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
	result, err := spec.ReadSpec(data)

	printSpecParseResult(result)
	fail.IfErrorF(err, "Failed to parse spec: %s", specFile)
	return result.Spec
}
