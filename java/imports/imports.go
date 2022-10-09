package imports

import (
	"generator"
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

func (self *imports) AddStatic(imports ...string) *imports {
	for _, i := range imports {
		self.imports = append(self.imports, "static "+i)
	}
	return self
}

func (self *imports) Write(w generator.Writer) {
	for _, imp := range self.imports {
		w.Line(`import %s;`, imp)
	}
}
