package cmd

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
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
	console.PrintLnF("Parsing spec file: %s", specFile)
	result, err := spec.ReadSpecFile(specFile)

	fail.IfErrorF(err, "Failed to read spec file: %s", specFile)
	printSpecParseResult(result)
	if err != nil {
		fail.FailF("Failed to parse spec: %s", specFile)
	}
	return result.Spec
}