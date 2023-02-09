import * as t from '../test-service/io-ts'

import * as assert from 'uvu/assert'

export const checkEncodeDecode = <A, O, I>(theType: t.Type<A, O, I>, decoded: A, encoded: I) => {
  let encodedActual = t.encode(theType, decoded)
  assert.equal(encoded, encodedActual)
  let decodedActual = t.decode(theType, encoded)
  assert.equal(decoded, decodedActual)
}

export function datetime(year: number, month: number, date: number, hours: number, minutes: number, seconds: number): Date {
  const utcDatetime = Date.UTC(year, month-1, date, hours, minutes, seconds)
  const result = new Date(utcDatetime) 
  return result
}