package iots

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/typescript/v2/modules"
)

var IoTs = "io-ts"

func (g *Generator) SetupLibrary(module modules.Module) *generator.CodeFile {
	code := `
export * from 'io-ts'
export * from 'io-ts-types'

import * as t from 'io-ts'

import { pipe } from 'fp-ts/lib/function'
import { map, chain, match } from 'fp-ts/lib/Either'

enum Enum {}

export class EnumType<E extends typeof Enum> extends t.Type<E[keyof E]> {
  readonly _tag: 'EnumType' = 'EnumType'
  private readonly _enum: E
  private readonly _enumValues: Set<string | number>
  constructor(e: E, name?: string) {
    super(
      name || 'enum',
      (u): u is E[keyof E] => {
        if (!this._enumValues.has(u as any)) return false
        // Don't allow key names from number enum reverse mapping
        if (typeof (this._enum as any)[u as string] === 'number') return false
        return true
      },
      (u, c) => (this.is(u) ? t.success(u) : t.failure(u, c)),
      t.identity
    )
    this._enum = e
    this._enumValues = new Set(Object.values(e))
  }
}

const enumType = <E extends typeof Enum>(e: E, name?: string) => new EnumType<E>(e, name)

export { enumType as enum }

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

export interface DateISOStringNoTimezoneC extends t.Type<Date, string, unknown> {}

export const DateISOStringNoTimezone: DateISOStringNoTimezoneC = new t.Type<Date, string, unknown>(
  'DateISOStringNoTimezone',
  (u): u is Date => u instanceof Date,
  (u, c) =>
    pipe(
      t.string.validate(u, c),
      chain(s => {
        const d = stringToDatetime(s)
        return isNaN(d.getTime()) ? t.failure(u, c) : t.success(d)
      })
    ),
  a => datetimeToString(a)
)

export class WithDefault<RT extends t.Any, A = any, O = A, I = unknown> extends t.Type<A, O, I> {
  readonly _tag: 'WithDefault' = 'WithDefault'
  constructor(
      name: string,
      is: WithDefault<RT, A, O, I>['is'],
      validate: WithDefault<RT, A, O, I>['validate'],
      serialize: WithDefault<RT, A, O, I>['encode'],
      readonly type: RT,
  ) {
      super(name, is, validate, serialize)
  }
}

export const withDefault = <RT extends t.Type<A, O>, A = any, O = A>(type: RT, defaultValue: t.TypeOf<RT>): WithDefault<RT, t.TypeOf<RT>, t.OutputOf<RT>, unknown> => {
  const Nullable = t.union([type, t.null, t.undefined])
  return new WithDefault(
    'WithDefault',
      (m: unknown): m is t.TypeOf<RT> => type.is(m),
      (s: unknown, c: t.Context): t.Validation<t.TypeOf<RT>> => {
          const validationResult: t.Validation<t.TypeOf<RT | t.NullC | t.UndefinedC>> = Nullable.validate(s, c)
          const applyDefault = map<A | null | undefined, A>(value => value != null ? value : defaultValue)
          return applyDefault(validationResult)
      },
      (a: t.TypeOf<RT>) => type.encode(a),
      type,
  )
}

export class DecodeError extends Error {
    errors: ValidationError[]
    constructor(errors: ValidationError[]) {
        super('Decoding failed')
        this.errors = errors
    }
}

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

const convertError = (error: t.ValidationError): ValidationError => {
  let code = "parsing_failed"
  if (error.value === undefined) {
    code = "missing"
  }
  const path = error.context[error.context.length-1].key
//  const path = error.context.slice(1).map(item => item.key).join(".")
  return {
    code,
    path,
    message: "Failed to parse"
  }
}

export const decodeR = <A, O, I>(codec: t.Type<A, O, I>, value: I): Result<A, ValidationError[]> => {
  return pipe(
    codec.decode(value),
    match(
      (errors: t.ValidationError[]): Result<A, ValidationError[]> => { return { error: errors.map(convertError) } },
      (value: A): Result<A, ValidationError[]> => { return { value } },
    )
  )
}

export const decode = <A, O, I>(codec: t.Type<A, O, I>, value: I): A => {
  const result = decodeR(codec, value)
  if (result.error) {
    throw new DecodeError(result.error)
  }
  return result.value
}

export const encode = <A, O, I>(codec: t.Type<A, O, I>, value: A): O => {
    return codec.encode(value)
}
`
	return &generator.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
