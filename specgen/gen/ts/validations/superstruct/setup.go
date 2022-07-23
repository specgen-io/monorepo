package superstruct

import (
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

var Superstruct = "superstruct"

func (g *Generator) SetupLibrary(module modules.Module) *generator.CodeFile {
	code := `
export * from "superstruct"

import * as t from "superstruct"

interface Success<T> {
    value: T
    error?: never
}

interface Error<E> {
    value?: never
    error: E
}

export type Result<T, E> = Success<T> | Error<E>

export interface ValidationError {
    path: string
    code: string
    message?: string
}

const convertFailure = (failure: t.Failure): ValidationError => {
    let code = "unknown"
    if (failure.message.includes("but received: undefined")) {
        code = "missing"
    } else if (failure.message.includes("but received:")) {
        code = "parsing_failed"
    }
    return {
        code, 
        path: failure.path.join("."), 
        message: failure.message
    }  
}

export const decodeR = <T>(struct: t.Struct<T, unknown>, value: unknown): Result<T, ValidationError[]> => {
    try {
        return {value: t.create(value, struct)}
    } catch (err: unknown) {
        const structError = err as t.StructError
        return {error: structError.failures().map(convertFailure)}
    }
}

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

	return &generator.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
