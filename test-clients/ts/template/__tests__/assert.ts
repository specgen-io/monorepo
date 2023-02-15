export * from 'uvu/assert'

import { AssertionError } from 'assert';
import * as assert from 'uvu/assert'

export const asyncThrows = async (theFunction: () => Promise<any>, errorType: any): Promise<any> => {
  try {
    await theFunction();
    assert.unreachable('does not throw');
  } catch (err) {
    if (err instanceof assert.Assertion) {
      throw err;
    }
    assert.instance(err, errorType)
    return err
  }
  throw new AssertionError({message: "never supposed to get here"})
}