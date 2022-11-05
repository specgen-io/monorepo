package writer

import (
	"fmt"
	"generator"
	"scala/packages"
)

func ScalaConfig() generator.Config {
	return generator.Config{"  ", 2, nil}
}

func New(thepackage packages.Package, className string) generator.Writer {
	config := ScalaConfig()
	filename := thepackage.GetPath(fmt.Sprintf("%s.scala", className))
	config.Substitutions["[[.ClassName]]"] = className
	w := generator.NewWriter2(filename, config)
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	return w
}
