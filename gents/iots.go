package gents

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

var IoTs = "io-ts"

var importIoTsEncoding = `import * as t from './io-ts'`

func generateIoTsStaticCode(path string) *gen.TextFile {
	code := `
export * from 'io-ts';
export * from 'io-ts-types';

import * as t from 'io-ts';

// Copy-paste from here: https://github.com/gcanti/io-ts/pull/366

enum Enum {}
/**
 * @since 2.3.0
 */
export class EnumType<E extends typeof Enum> extends t.Type<E[keyof E]> {
  /**
   * @since 2.3.0
   */
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

/**
 * @since 2.3.0
 */
const enumType = <E extends typeof Enum>(e: E, name?: string) => new EnumType<E>(e, name);

export { enumType as enum }

import { pipe } from 'fp-ts/lib/pipeable';
import { fold } from 'fp-ts/lib/Either';
import { identity } from 'fp-ts/lib/function';

export class DecodeError extends Error {
    errors: t.Errors
    constructor(errors: t.Errors) {
        super('Decoding failed');
        this.errors = errors;
    }
}

export const decode = <A, O, I>(codec: t.Type<A, O, I>, value: I): A => {
    return pipe(
        codec.decode(value),
        fold(
            errors => { throw new DecodeError(errors); },
            identity
        )
    );
};

export const encode = <A, O, I>(codec: t.Type<A, O, I>, value: A): O => {
    return codec.encode(value);
};
`
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
