import * as t from '../test-service/io-ts'
import { checkEncodeDecode } from './util'

import {
  Message, TMessage
} from '../test-service/v2/models';

describe('v2 object', function() {
  let decoded: Message = {field: 'the string'}
  let encoded = {'field': 'the string'}
  checkEncodeDecode(TMessage, decoded, encoded)
});