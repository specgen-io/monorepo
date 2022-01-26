package imports

import (
	"github.com/specgen-io/specgen/v2/sources"
)

type imports struct {
	imports []string
}

func New() *imports {
	return &imports{imports: []string{}}
}

func (self *imports) Add(imports ...string) *imports {
	self.imports = append(self.imports, imports...)
	return self
}

func (self *imports) Write(w *sources.Writer) {
	for _, imp := range self.imports {
		w.Line(`import %s;`, imp)
	}
}
