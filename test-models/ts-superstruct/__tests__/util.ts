import * as t from '../test-service/superstruct'

export const checkEncodeDecode = <T>(theType: t.Struct<T, unknown>, decoded: T, encoded: unknown) => {
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