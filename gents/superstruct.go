package gents

import (
	"gopkg.in/specgen-io/specgen.v2/gen"
	"strings"
)

func generateSuperstruct(path string) *gen.TextFile {
	code := `
export * from 'superstruct'

import * as t from 'superstruct';

export const DateTime = t.coerce<Date, unknown, string>(t.date(), t.string(), (value: string): Date => new Date(value))

export const decode = <T>(struct: t.Struct<T, unknown>, value: unknown): T => {
    return t.create(value, struct)
};

export const encode = <T>(struct: t.Struct<T, unknown>, value: T): unknown => {
    return JSON.parse(JSON.stringify(value))
};
`
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
