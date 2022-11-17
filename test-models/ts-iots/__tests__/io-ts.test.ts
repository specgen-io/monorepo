import * as t from '../test-service/io-ts'
import { checkEncodeDecode } from './util'

const TDefaultedField = t.type({
    theDefaulted: t.withDefault(t.string, "the default value"),
})

type DefaultedField = t.TypeOf<typeof TDefaultedField>

describe('withDefault', function() {
    it('decode absent value', function() {
        let decoded: DefaultedField = {theDefaulted: 'the default value'}
        let encoded = {}
        expect(t.decode(TDefaultedField, encoded)).toStrictEqual(decoded)
    })

    it('decode undefined value', function() {
        let decoded: DefaultedField = {theDefaulted: 'the default value'}
        let encoded = {theDefaulted: undefined}
        expect(t.decode(TDefaultedField, encoded)).toStrictEqual(decoded)
    })
    
    it('decode null value', function() {
        let decoded: DefaultedField = {theDefaulted: 'the default value'}
        let encoded = {theDefaulted: null}
        expect(t.decode(TDefaultedField, encoded)).toStrictEqual(decoded)
    })

    it('decode present value', function() {
        let decoded: DefaultedField = {theDefaulted: "there's some value"}
        let encoded = {theDefaulted: "there's some value"}
        expect(t.decode(TDefaultedField, encoded)).toStrictEqual(decoded)
    })

    it('decode fail', function() {
        expect(() => t.decode(TDefaultedField, {theDefaulted: 123})).toThrowError('Decoding failed');
      })
    
})