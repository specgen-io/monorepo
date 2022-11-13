import { test } from 'uvu';
import * as assert from 'uvu/assert';

import {Message, Choice, Parameters, UrlParameters} from '../test-service/models'
import {client as echoClient} from '../test-service/echo'
import {client as checkClient} from '../test-service/check'

const config = {baseURL: process.env.SERVICE_URL!}

test('echoBodyString', async function() {
  let body: string = "some text"
  const client = echoClient(config)
  let response = await client.echoBodyString({body})
  assert.equal(response, body, 'response matches request')
})

test('echoBodyModel', async function() {
  const client = echoClient(config)
  let body: Message = {int_field: 123, string_field: "the string"}
  let response = await client.echoBodyModel({body})
  assert.equal(response, body, 'response matches request')
})

test('echoBodyArray', async function() {
  const client = echoClient(config)
  let body: string[] = ["the str1", "the str1"]
  let response = await client.echoBodyArray({body})
  assert.equal(response, body, 'response matches request')
})

test('echoBodyMap', async function() {
  const client = echoClient(config)
  let body: Record<string, string> = {"one": "the str1", "two": "the str1"}
  let response = await client.echoBodyMap({body})
  assert.equal(response, body, 'response matches request')
})

test('echoQuery', async function() {
  const client = echoClient(config)
  let expected: Parameters = {
    int_field: 123, 
    long_field: 12345,
    float_field: 1.23,
    double_field: 12.345,
    decimal_field: 12345,
    bool_field: true,
    string_field: "the value",
    string_opt_field: "the value",
    string_defaulted_field: "value",
    string_array_field: ["the str1", "the str2"],
    uuid_field: "123e4567-e89b-12d3-a456-426655440000",
    date_field: "2021-01-01",
    date_array_field: ["2021-01-02"],
    datetime_field: new Date("2021-01-02T23:54"),
    enum_field: Choice.SECOND_CHOICE,
  }
  let response = await client.echoQuery({
    intQuery: 123, 
    longQuery: 12345,
    floatQuery: 1.23,
    doubleQuery: 12.345,
    decimalQuery: 12345,
    boolQuery: true,
    stringQuery: "the value",
    stringOptQuery: "the value",
    stringDefaultedQuery: "value",
    stringArrayQuery: ["the str1", "the str2"],
    uuidQuery: "123e4567-e89b-12d3-a456-426655440000",
    dateQuery: "2021-01-01",
    dateArrayQuery: ["2021-01-02"],
    datetimeQuery: new Date("2021-01-02T23:54"),
    enumQuery: Choice.SECOND_CHOICE,
  })
  assert.equal(response, expected, 'response matches expected')
})

test('echoHeader', async function() {
  const client = echoClient(config)
  let expected: Parameters = {
    int_field: 123, 
    long_field: 12345,
    float_field: 1.23,
    double_field: 12.345,
    decimal_field: 12345,
    bool_field: true,
    string_field: "the value",
    string_opt_field: "the value",
    string_defaulted_field: "value",
    string_array_field: ["the str1", "the str2"],
    uuid_field: "123e4567-e89b-12d3-a456-426655440000",
    date_field: "2021-01-01",
    date_array_field: ["2021-01-02"],
    datetime_field: new Date("2021-01-02T23:54"),
    enum_field: Choice.SECOND_CHOICE, 
  }
  let response = await client.echoHeader({
    intHeader: 123, 
    longHeader: 12345,
    floatHeader: 1.23,
    doubleHeader: 12.345,
    decimalHeader: 12345,
    boolHeader: true,
    stringHeader: "the value",
    stringOptHeader: "the value",
    stringDefaultedHeader: "value",
    stringArrayHeader: ["the str1", "the str2"],
    uuidHeader: "123e4567-e89b-12d3-a456-426655440000",
    dateHeader: "2021-01-01",
    dateArrayHeader: ["2021-01-02"],
    datetimeHeader: new Date("2021-01-02T23:54"),
    enumHeader: Choice.SECOND_CHOICE,
  })
  assert.equal(response, expected, 'response matches expected')
})

test('echoUrlParams', async function() {
  const client = echoClient(config)
  let expected: UrlParameters = {
    int_field: 123, 
    long_field: 12345,
    float_field: 1.23,
    double_field: 12.345,
    decimal_field: 12345,
    bool_field: true,
    string_field: "the value",
    uuid_field: "123e4567-e89b-12d3-a456-426655440000",
    date_field: "2021-01-01",
    datetime_field: new Date("2021-01-02T23:54"),
    enum_field: Choice.SECOND_CHOICE,
  }
  let response = await client.echoUrlParams({
    intUrl: 123, 
    longUrl: 12345,
    floatUrl: 1.23,
    doubleUrl: 12.345,
    decimalUrl: 12345,
    boolUrl: true,
    stringUrl: "the value",
    uuidUrl: "123e4567-e89b-12d3-a456-426655440000",
    dateUrl: "2021-01-01",
    datetimeUrl: new Date("2021-01-02T23:54"),
    enumUrl: Choice.SECOND_CHOICE,
  })
  assert.equal(response, expected, 'response matches expected')
})

test('checkEmpty', async function() {
  let client = checkClient(config)
  let response = await client.checkEmpty()
  assert.is(response, undefined, 'response on check empty is void')
})

test.run();