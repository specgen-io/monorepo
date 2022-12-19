import {client} from '../test-service/check'
import {BadRequestException, ConflictException, ForbiddenException} from '../test-service/errors'

import { test } from 'uvu'
import * as assert from './assert'

const checkClient = client({baseURL: process.env.SERVICE_URL!, timeout: 20000})

test('checkEmpty', async function() {
  assert.not.throws(async () => {
    await checkClient.checkEmpty()
  })
})

test('checkEmptyResponse', async function() {
  assert.not.throws(async () => {
    await checkClient.checkEmptyResponse({body: {int_field: 123, string_field: "the string"}})
  })
})


test('checkBadRequest', async function() {
  const err = await assert.asyncThrows(checkClient.checkBadRequest, BadRequestException) as BadRequestException
  assert.equal(err.message, 'Error response with status code 400')
})

test('checkForbidden', async function() {
  const err = await assert.asyncThrows(checkClient.checkForbidden, ForbiddenException)
})

test('checkConflict', async function() {
  const err = await assert.asyncThrows(checkClient.checkConflict, ConflictException)
})