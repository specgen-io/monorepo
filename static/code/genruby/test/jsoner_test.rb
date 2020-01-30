require "test/unit"

require 'dataclass'
require 'jsoner'

class PlainTypesDeserialization < Test::Unit::TestCase
  def test_deserialize_integer
    data = Jsoner.from_json(Integer, '123')
    T.check(Integer, data)
    assert_equal 123, data, "Integer should be parsable from JSON"
  end

  def test_deserialize_integer_fail
    assert_raise JsonerError do
      Jsoner.from_json(Integer, '"abc"')
    end
  end

  def test_deserialize_float
    data = Jsoner.from_json(Float, '1.23')
    T.check(Float, data)
    assert_equal 1.23, data, "Float should be parsable from JSON"
  end

  def test_deserialize_float_from_integer
    data = Jsoner.from_json(Float, '123')
    T.check(Float, data)
    assert_equal 123.0, data, "Float should be parsable from JSON"
  end

  def test_deserialize_float_fail
    assert_raise JsonerError do
      Jsoner.from_json(Float, '"abc"')
    end
  end

  def test_deserialize_boolean
    data = Jsoner.from_json(Boolean, 'true')
    T.check(Boolean, data)
    assert_equal true, data, "Boolean should be parsable from JSON"
  end

  def test_deserialize_boolean_fail
    assert_raise JsonerError do
      Jsoner.from_json(Boolean, '1')
    end
  end

  def test_deserialize_string
    data = Jsoner.from_json(String, '"the string"')
    T.check(String, data)
    assert_equal "the string", data, "String should be parsable from JSON"
  end

  def test_deserialize_string_fail
    assert_raise JsonerError do
      Jsoner.from_json(String, '123')
    end
  end

  def test_deserialize_uuid
    data = Jsoner.from_json(UUID, '"58d5e212-165b-4ca0-909b-c86b9cee0111"')
    T.check(UUID, data)
    assert_equal "58d5e212-165b-4ca0-909b-c86b9cee0111", data, "String should be parsable from JSON"
  end

  def test_deserialize_uuid_fail
    assert_raise JsonerError do
      Jsoner.from_json(UUID, '"abc"')
    end
  end

  def test_deserialize_date
    data = Jsoner.from_json(Date, '"2019-11-30"')
    T.check(Date, data)
    assert_equal Date.new(2019, 11, 30), data, "Date should be parsable from JSON"
  end

  def test_deserialize_date_wrong_format
    assert_raise JsonerError do
      Jsoner.from_json(Date, '"11/30/2019"')
    end
  end

  def test_deserialize_datetime
    data = Jsoner.from_json(DateTime, '"2019-11-30T17:45:55+00:00"')
    T.check(DateTime, data)
    assert_equal DateTime.new(2019, 11, 30, 17, 45, 55), data, "DateTime should be parsable from JSON"
  end

  def test_deserialize_datetime_wrong_format
    assert_raise JsonerError do
      Jsoner.from_json(DateTime, '"2019-11-30 17:45:55"')
    end
  end

  def test_deserialize_nilable_nil
    data = Jsoner.from_json(T.nilable(String), 'null')
    T.check(T.nilable(String), data)
    assert_equal nil, data, "Nilable null value should be parsable from JSON"
  end

  def test_deserialize_nilable_string
    data = Jsoner.from_json(T.nilable(String), '"some string"')
    T.check(T.nilable(String), data)
    assert_equal "some string", data, "Nilable non-null value should be parsable from JSON"
  end

  def test_deserialize_non_nilable
    assert_raise JsonerError do
      Jsoner.from_json(String, 'null')
    end
  end
end

class PlainTypesSerialization < Test::Unit::TestCase
  def test_serialize_integer
    assert_equal "123", Jsoner.to_json(Integer,123), "Integer should be serializable to JSON"
  end

  def test_serialize_integer_fail
    assert_raise JsonerError do
      Jsoner.to_json(Integer,123.4)
    end
  end

  def test_serialize_float
    assert_equal "12.3", Jsoner.to_json(Float,12.3), "Float should be serializable to JSON"
  end

  def test_serialize_float_fail
    assert_raise JsonerError do
      Jsoner.to_json(Float,123)
    end
  end

  def test_serialize_boolean
    assert_equal "true", Jsoner.to_json(Boolean, true), "Boolean should be serializable to JSON"
  end

  def test_serialize_string
    assert_equal '"the string"', Jsoner.to_json(String,"the string"), "String should be serializable to JSON"
  end

  def test_serialize_date
    assert_equal '"2019-11-30"', Jsoner.to_json(Date, Date.new(2019, 11, 30)), "Date should be serializable to JSON"
  end

  def test_serialize_date_fail
    assert_raise JsonerError do
      Jsoner.to_json(Date, 123)
    end
  end

  def test_serialize_datetime
    assert_equal '"2019-11-30T17:45:55"', Jsoner.to_json(DateTime, DateTime.new(2019, 11, 30, 17, 45, 55)), "DateTime should be serializable to JSON"
  end

  def test_serialize_nil
    json_str = Jsoner.to_json(T.nilable(Untyped), nil)
    assert_equal "null", json_str, "nil should be serializable to JSON"
  end

  def test_serialize_uuid
    assert_equal '"58d5e212-165b-4ca0-909b-c86b9cee0111"', Jsoner.to_json(UUID, "58d5e212-165b-4ca0-909b-c86b9cee0111"), "UUID should be serializable to JSON"
  end
end

class ArrayDeserialization < Test::Unit::TestCase
  def test_deserialize_array
    data = Jsoner.from_json(T.array(String), '["the string", "the other string"]')
    T.check(T.array(String), data)
    assert_equal ["the string", "the other string"], data, "Should parse array of strings"
  end

  def test_deserialize_array_of_datetime
    data = Jsoner.from_json(T.array(DateTime), '["2019-11-30T17:45:55+00:00"]')
    T.check(T.array(DateTime), data)
    assert_equal [DateTime.new(2019, 11, 30, 17, 45, 55)], data, "Should parse array of DateTime"
  end
end

class ArraySerialization < Test::Unit::TestCase
  def test_serialize_array
    assert_equal '["the string","the other string"]', Jsoner.to_json(T.array(String), ["the string", "the other string"]), "Array should be serializable to JSON"
  end

  def test_serialize_array_of_datetime
    assert_equal '["2019-11-30T17:45:55"]', Jsoner.to_json(T.array(DateTime), [DateTime.new(2019, 11, 30, 17, 45, 55)]), "Array should be serializable to JSON"
  end
end

class HashDeserialization < Test::Unit::TestCase
  def test_deserialize_hash
    data = Jsoner.from_json(T.hash(String, Integer), '{"one": 123, "two": 456}')
    T.check(T.hash(String, Integer), data)
    assert_equal ({"one" => 123, "two" => 456}), data, "Should parse hash"
  end

  def test_deserialize_hash_string_to_datetime
    data = Jsoner.from_json(T.hash(String, DateTime), '{"one": "2019-11-30T17:45:55+00:00"}')
    T.check(T.hash(String, DateTime), data)
    assert_equal ({"one" => DateTime.new(2019, 11, 30, 17, 45, 55)}), data, "Should parse hash"
  end
end

class HashSerialization < Test::Unit::TestCase
  def test_serialize_hash
    assert_equal '{"one":123,"two":456}', Jsoner.to_json(T.hash(String, Integer), {"one" => 123, "two" => 456}), "Hash should be serializable to JSON"
  end

  def test_serialize_hash_string_to_datetime
    assert_equal '{"one":"2019-11-30T17:45:55"}', Jsoner.to_json(T.hash(String, DateTime), {"one" => DateTime.new(2019, 11, 30, 17, 45, 55)}), "Hash should be serializable to JSON"
  end
end

class AnyDeserialization < Test::Unit::TestCase
  def test_deserialize_any
    data = Jsoner.from_json(T.any(String, Integer), '"the string"')
    T.check(T.any(String, Integer), data)
    assert_equal "the string", data, "Should parse any type"
  end

  def test_deserialize_any_fail
    assert_raise JsonerError do
      Jsoner.from_json(T.any(Integer, Float), '"the string"')
    end
  end
end

class AnySerialization < Test::Unit::TestCase
  def test_serialize_any
    data = Jsoner.to_json(T.any(String, Integer), "the string")
    T.check(T.any(String, Integer), data)
    assert_equal '"the string"', data, "Should serialize any type"
  end
end