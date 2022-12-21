package writer

import (
	"fmt"
	"strings"
	"typescript/module"
)

type imports struct {
	stars         map[string]module.Module
	starsOrder    []string
	names         map[string][]string
	namesOrder    []string
	defaults      map[string]string
	defaultsOrder []string
	Target        module.Module
}

func NewImports(target module.Module) *imports {
	return &imports{map[string]module.Module{}, []string{}, map[string][]string{}, []string{}, map[string]string{}, []string{}, target}
}

func (self *imports) Star(m module.Module, alias string) {
	self.stars[alias] = m
	self.starsOrder = append(self.starsOrder, alias)
}

func (self *imports) Names(m module.Module, names ...string) {
	self.LibNames(m.GetImport(self.Target), names...)
}

func (self *imports) Aliased(m module.Module, name, alias string) {
	self.LibNames(m.GetImport(self.Target), fmt.Sprintf(`%s as %s`, name, alias))
}

func (self *imports) LibNames(lib string, names ...string) {
	moduleNames, found := self.names[lib]
	if !found {
		self.namesOrder = append(self.namesOrder, lib)
		moduleNames = []string{}
	}
	self.names[lib] = append(moduleNames, names...)
}

func (self *imports) Default(m string, name string) {
	self.defaults[m] = name
	self.defaultsOrder = append(self.defaultsOrder, m)
}

func (self *imports) Lines() []string {
	lines := []string{}
	for _, m := range self.defaultsOrder {
		lines = append(lines, fmt.Sprintf(`import %s from '%s'`, self.defaults[m], m))
	}
	for _, alias := range self.starsOrder {
		lines = append(lines, fmt.Sprintf("import * as %s from '%s'", alias, self.stars[alias].GetImport(self.Target)))
	}
	for _, m := range self.namesOrder {
		if len(self.names[m]) > 0 {
			lines = append(lines, fmt.Sprintf(`import { %s } from '%s'`, strings.Join(self.names[m], ", "), m))
		}
	}
	return lines
}
