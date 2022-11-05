package imports

import (
	"fmt"
	"generator"
	"strings"
	"typescript/module"
)

type imports struct {
	stars    map[string]module.Module
	names    map[string][]string
	defaults map[string]string
	Target   module.Module
}

func New(target module.Module) *imports {
	return &imports{map[string]module.Module{}, map[string][]string{}, map[string]string{}, target}
}

func (self *imports) Star(m module.Module, alias string) {
	self.stars[alias] = m
}

func (self *imports) Names(m module.Module, names ...string) {
	self.LibNames(m.GetImport(self.Target), names...)
}

func (self *imports) Aliased(m module.Module, name, alias string) {
	self.LibNames(m.GetImport(self.Target), fmt.Sprintf(`%s as %s`, name, alias))
}

func (self *imports) LibNames(lib string, names ...string) {
	moduleNames, found := self.names[lib]
	if found {
		moduleNames = []string{}
	}
	self.names[lib] = append(moduleNames, names...)
}

func (self *imports) Default(m string, name string) {
	self.defaults[m] = name
}

func (self *imports) Write(w generator.Writer) {
	for alias, m := range self.stars {
		w.Line("import * as %s from '%s'", alias, m.GetImport(self.Target))
	}
	for m, names := range self.names {
		if len(names) > 0 {
			w.Line(`import { %s } from '%s'`, strings.Join(names, ", "), m)
		}
	}
	for m, name := range self.defaults {
		w.Line(`import %s from '%s'`, name, m)
	}
}
