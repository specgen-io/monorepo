package gen

import (
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/openapi/v2"
	"github.com/specgen-io/specgen/ruby/v2"
	"github.com/specgen-io/specgen/scala/v2"
	"github.com/specgen-io/specgen/v2/gen/golang"
	"github.com/specgen-io/specgen/v2/gen/java"
	"github.com/specgen-io/specgen/v2/gen/kotlin"
	"github.com/specgen-io/specgen/v2/gen/ts"
)

var Generators = []generator.Generator{
	golang.Models,
	golang.Client,
	golang.Service,
	java.Models,
	java.Client,
	java.Service,
	kotlin.Models,
	kotlin.Client,
	kotlin.Service,
	ruby.Models,
	ruby.Client,
	scala.Models,
	scala.Client,
	scala.Service,
	ts.Models,
	ts.Client,
	ts.Service,
	openapi.Openapi,
}
