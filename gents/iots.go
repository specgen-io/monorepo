package gents

import (
	"specgen/gen"
	"strings"
)

func generateIoTs(path string) *gen.TextFile {
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
`
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
