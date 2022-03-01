package service

import (
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/types"
)

var Spring = "spring"

type SpringGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewSpringGenerator(types *types.Types, models models.Generator) *SpringGenerator {
	return &SpringGenerator{types, models}
}
