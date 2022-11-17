import * as t from '../test-service/io-ts'

export const checkEncodeDecode = <A, O, I>(theType: t.Type<A, O, I>, decoded: A, encoded: I) => {
    it('encode', function() {
    expect(t.encode(theType, decoded)).toStrictEqual(encoded);
  })
  it('decode', function() {
    expect(t.decode(theType, encoded)).toStrictEqual(decoded);
  })
}

export function datetime(year: number, month: number, date: number, hours: number, minutes: number, seconds: number): Date {
  const utcDatetime = Date.UTC(year, month-1, date, hours, minutes, seconds)
  const result = new Date(utcDatetime) 
  return result
}