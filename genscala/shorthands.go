package genscala

import (
	"github.com/vsapronov/gopoetry/scala"
	"strings"
)

var Code = scala.Code
var Line = scala.Line
var Block = scala.Block
var Scope = scala.Scope
var If = scala.If
var Statements = scala.Statements
var Class = scala.Class
var CaseClass = scala.CaseClass
var Object = scala.Object
var CaseObject = scala.CaseObject
var Trait = scala.Trait
var Param = scala.Param
var Import = scala.Import
var Def = scala.Def
var Constructor = scala.Constructor
var Unit = scala.Unit
var Val = scala.Val
var Lazy = scala.Lazy
var NoCode = scala.NoCode
var Attribute = scala.Attribute

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}