import { checkEncodeDecode } from './util'
import { test } from 'uvu'

import { Message, TMessage } from '../src/spec/v2/models';

test('v2 object encode + decode', function() {
  let decoded: Message = {field: 'the string'}
  let encoded = {'field': 'the string'}
  checkEncodeDecode(TMessage, decoded, encoded)
})