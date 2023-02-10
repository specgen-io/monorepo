import { checkEncodeDecode } from './util'

import { test } from 'uvu'

import {
  OptionalFields,
  TOptionalFields,
} from '../src/spec/models';


test('optional fields: encode + decode undefined values', function() {
  let decoded: OptionalFields = {int_option_field: undefined, string_option_field: undefined}
  let encoded = {int_option_field: undefined, string_option_field: undefined}
  checkEncodeDecode(TOptionalFields, decoded, encoded)
})

test.run()