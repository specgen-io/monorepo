package gents

import (
	"specgen/gen"
	"strings"
)

func generateCodec(path string) *gen.TextFile {
	code := `
import { Errors, Type } from 'io-ts';
import { pipe } from 'fp-ts/lib/pipeable';
import { fold } from 'fp-ts/lib/Either';
import { identity } from 'fp-ts/lib/function';

export class DecodeError extends Error {
    errors: Errors
    constructor(errors: Errors) {
        super('Decoding failed');
        this.errors = errors;
    }
}

export const decode = <A, O, I>(codec: Type<A, O, I>, value: I): A => {
    return pipe(
        codec.decode(value),
        fold(
            errors => { throw new DecodeError(errors); },
            identity
        )
    );
};

export const encode = <A, O, I>(codec: Type<A, O, I>, value: A): O => {
    return codec.encode(value);
};
`

	return &gen.TextFile{path, strings.TrimSpace(code)}
}
