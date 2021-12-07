package conoldspec

import (
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/spec/old"
)

func ConvertFromOldSpec(inSpecFile, outSpecFile string) error {
	return old.FormatSpec(inSpecFile, outSpecFile, spec.SpecVersion)
}
