package writer

import (
	"fmt"
	"strings"
	"typescript/module"
)

type imports struct {
	stars    map[string]module.Module
	names    map[string][]string
	defaults map[string]string
	Target   module.Module
}

func NewImports(target module.Module) *imports {
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

func (self *imports) Lines() []string {
	lines := []string{}
	for alias, m := range self.stars {
		lines = append(lines, fmt.Sprintf("import * as %s from '%s'", alias, m.GetImport(self.Target)))
	}
	for m, names := range self.names {
		if len(names) > 0 {
			lines = append(lines, fmt.Sprintf(`import { %s } from '%s'`, strings.Join(names, ", "), m))
		}
	}
	for m, name := range self.defaults {
		lines = append(lines, fmt.Sprintf(`import %s from '%s'`, name, m))
	}
	return lines
}
