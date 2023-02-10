import * as t from '../src/spec/{{validation.value}}'
import { checkEncodeDecode, datetime } from './util'

import { test } from 'uvu'
import * as assert from 'uvu/assert'

import {
  TMessage,
  Message,
  Choice,
  EnumFields,
  TEnumFields,
  Parent,
  TParent,
  NumericFields,
  TNumericFields,
  NonNumericFields,
  TNonNumericFields,
  ArrayFields,
  TArrayFields,
  MapFields,
  TMapFields,
  OptionalFields,
  TOptionalFields,
  OrderEventWrapper,
  TOrderEventWrapper,
  OrderEventDiscriminator,
  TOrderEventDiscriminator
} from '../src/spec/models';


test('object: encode + decode', function() {
  let decoded: Message = {field: 123}
  let encoded = {'field': 123}
  checkEncodeDecode(TMessage, decoded, encoded)
})

test('object: decode missing required fields', function() {
  assert.throws(() => t.decode(TMessage, {}))
})

test('enum fields: encode + decode', function() {
  let decoded: EnumFields = {enum_field: Choice.THIRD_CHOICE}
  let encoded = {'enum_field': 'Three'}
  checkEncodeDecode(TEnumFields, decoded, encoded)
})

test('enum fields: decode no required fields', function() {
  assert.throws(() => t.decode(TEnumFields, {}))
})

test('nested object: encode + decode', function() {
  let decoded: Parent = {'field': 'the string', 'nested': {'field': 123}}
  let encoded = {'field': 'the string', 'nested': {'field': 123}}
  checkEncodeDecode(TParent, decoded, encoded)
})

test('nested object: required is missing', function() {
  assert.throws(() => t.decode(TParent, {'field': 'the string'}))
})

test('numeric fields: encode + decode', function() {
  let decoded: NumericFields = {
    int_field: 123,
    long_field: 123,
    float_field: 12.3,
    double_field: 12.3,
    decimal_field: 12.3,
  }
  let encoded = {
    int_field: 123,
    long_field: 123,
    float_field: 12.3,
    double_field: 12.3,
    decimal_field: 12.3,
  }
  checkEncodeDecode(TNumericFields, decoded, encoded)
})

test('non numeric fields: encode + decode', function() {
  let decoded: NonNumericFields = {
    boolean_field: true,
    string_field: 'some string',
    uuid_field: '123e4567-e89b-12d3-a456-426655440000',
    date_field: '2021-01-01',
    datetime_field: datetime(2021, 1, 2, 4, 54, 0),
  }
  let encoded = {
    boolean_field: true,
    string_field: 'some string',
    uuid_field: '123e4567-e89b-12d3-a456-426655440000',
    date_field: '2021-01-01',
    datetime_field: '2021-01-02T04:54:00',
  }
  checkEncodeDecode(TNonNumericFields, decoded, encoded)
})

test('array fields: encode + decode', function() {
  let decoded: ArrayFields = {
    int_array_field: [1, 2, 3],
    string_array_field: ['one', 'two', 'three'],
  }
  let encoded = {
    int_array_field: [1, 2, 3],
    string_array_field: ['one', 'two', 'three'],
  }
  checkEncodeDecode(TArrayFields, decoded, encoded)
})

test('map fields: encode + decode', function() {
  let decoded: MapFields = {
    int_map_field: {'one': 1, 'two': 2, 'three': 3},
    string_map_field: {'one': 'first', 'two': 'second', 'three': 'third'},
  }
  let encoded = {
    int_map_field: {'one': 1, 'two': 2, 'three': 3},
    string_map_field: {'one': 'first', 'two': 'second', 'three': 'third'},
  }
  checkEncodeDecode(TMapFields, decoded, encoded)
})

test('optional fields: encode + decode', function() {
  let decoded: OptionalFields = {
    int_option_field: 123,
    string_option_field: 'the string',
  }
  let encoded = {
    int_option_field: 123,
    string_option_field: 'the string',
  }
  checkEncodeDecode(TOptionalFields, decoded, encoded)
})

test('optional fields: encode + decode null values', function() {
  let decoded: OptionalFields = {
    int_option_field: null,
    string_option_field: null,
  }
  let encoded = {
    int_option_field: null,
    string_option_field: null,
  }
  checkEncodeDecode(TOptionalFields, decoded, encoded)
})

test('oneof wrapped: encode + decode', function() {
  let decoded: OrderEventWrapper = { changed: { id: 'id123', quantity: 123 } }
  let encoded = { changed: {id: 'id123', quantity: 123} }
  checkEncodeDecode(TOrderEventWrapper, decoded, encoded)
})

test('oneof wrapped: missing all fields', function() {
  assert.throws(() => t.decode(TOrderEventWrapper, { created: {} }))
})

test('oneof discriminator: encode + decode', function() {
  let decoded: OrderEventDiscriminator = { _type: 'changed', id: 'id123', quantity: 123 }
  let encoded = { _type: 'changed', id: 'id123', quantity: 123 }
  checkEncodeDecode(TOrderEventDiscriminator, decoded, encoded)
})

test('oneof discriminator: missing all fields', function() {
  assert.throws(() => t.decode(TOrderEventDiscriminator, { _type: 'changed' }))
})

test.run()