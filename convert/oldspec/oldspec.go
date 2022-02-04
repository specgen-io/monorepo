package oldspec

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/spec/old"
	"io/ioutil"
)

func ConvertFromOldSpec(inSpecFile, outSpecFile string) error {
	inData, err := ioutil.ReadFile(inSpecFile)
	fail.IfErrorF(err, "Failed to read spec file: %s", inSpecFile)

	outData, err := old.FormatSpec(inData, spec.SpecVersion)
	fail.IfErrorF(err, "Failed to format spec: %s", inSpecFile)

	err = ioutil.WriteFile(outSpecFile, outData, 0644)
	fail.IfErrorF(err, "Failed to write spec to file: %s", outSpecFile)

	return nil
}
