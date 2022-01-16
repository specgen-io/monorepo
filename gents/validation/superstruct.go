package validation

import (
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

var Superstruct = "superstruct"

func generateSuperstructStaticCode(module modules.Module) *sources.CodeFile {
	code := `
export * from "superstruct"

import * as t from "superstruct"

export const decode = <T>(struct: t.Struct<T, unknown>, value: unknown): T => {
    return t.create(value, struct)
}

function datetimeToString(datetime: Date): string {
    return datetime.toISOString().slice(0, -5) //removing ending postfix .000Z
}

function stringToDatetime(str: string): Date {
    if (str.endsWith("Z")) {
        return new Date(str)
    } else {
        return new Date(str+"Z") // adding UTC symbol
    }
}

export const encode = <T>(_struct: t.Struct<T, unknown>, value: T): unknown => {
    return JSON.parse(JSON.stringify(value, function (this, key, value) {
        if (this[key] instanceof Date) {
            return datetimeToString(this[key])
        } else {
            return value
        }
    }))
}

export const StrDateTime = t.coerce<Date, unknown, string>(t.date(), t.string(), (value: string): Date => stringToDatetime(value))
export const StrInteger = t.coerce<number, unknown, string>(t.integer(), t.string(), (value: string): number => parseInt(value))
export const StrFloat = t.coerce<number, unknown, string>(t.number(), t.string(), (value: string): number => parseFloat(value))
export const StrBoolean = t.coerce<boolean, unknown, string>(t.boolean(), t.string(), (value: string): boolean => {
    switch (value) {
        case 'true': return true
        case 'false': return false
        default: throw new Error('Unknown boolean value: "'+value+'"')
    }
})`

	return &sources.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
