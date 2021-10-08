package gents

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

var Superstruct = "superstruct"

func generateSuperstructStaticCode(module module) *gen.TextFile {
	code := `
export * from "superstruct"

import * as t from "superstruct"

export const decode = <T>(struct: t.Struct<T, unknown>, value: unknown): T => {
    return t.create(value, struct)
}

export const encode = <T>(struct: t.Struct<T, unknown>, value: T): unknown => {
    return JSON.parse(JSON.stringify(value))
}

export const StrDateTime = t.coerce<Date, unknown, string>(t.date(), t.string(), (value: string): Date => new Date(value))
export const StrInteger = t.coerce<number, unknown, string>(t.integer(), t.string(), (value: string): number => parseInt(value))
export const StrFloat = t.coerce<number, unknown, string>(t.number(), t.string(), (value: string): number => parseFloat(value))
export const StrBoolean = t.coerce<boolean, unknown, string>(t.boolean(), t.string(), (value: string): boolean => {
    switch (value) {
        case 'true': return true
        case 'false': return false
        default: throw new Error('Unknown boolean value: "'+value+'"')
    }
})`

	return &gen.TextFile{module.GetPath(), strings.TrimSpace(code)}
}
