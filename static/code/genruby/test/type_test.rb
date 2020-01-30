require "test/unit"
require "date"

require 'type'
require 'enum'

class TypeEquality < Test::Unit::TestCase
  def test_plain_equals
    assert_true Integer == Integer
  end

  def test_plain_not_equals
    assert_false Integer == String
  end

  def test_uuid_equals
    assert_true UUID == UUID
  end

  def test_boolean_equals
    assert_true Boolean == Boolean
  end

  def test_untyped_equals
    assert_true Untyped == Untyped
  end

  def test_nilable_equals
    assert_true T.nilable(Integer) == T.nilable(Integer)
  end

  def test_nilable_not_equals
    assert_false T.nilable(Integer) == T.nilable(String)
  end

  def test_array_equals
    assert_true T.array(Integer) == T.array(Integer)
  end

  def test_array_not_equals
    assert_false T.array(Integer) == Integer
  end

  def test_array_other_item_type
    assert_false T.array(Integer) == T.array(String)
  end

  def test_hash_equals
    assert_true T.hash(String, Integer) == T.hash(String, Integer)
  end

  def test_hash_not_equals
    assert_false T.hash(String, Integer) == String
  end

  def test_hash_other_value_type
    assert_false T.hash(String, Integer) == T.hash(String, String)
  end

  def test_any_equals
    assert_true T.any(String, Integer) == T.any(String, Integer)
  end

  def test_any_equals_other_order
    assert_true T.any(String, Integer) == T.any(Integer, String)
  end

  def test_any_of_other_type
    assert_false T.any(String, Integer) == T.any(Integer, Float)
  end
end

class TypeCheck < Test::Unit::TestCase
  def test_nil
    assert_equal nil, T.check(NilClass, nil)
  end

  def test_nil_string
    assert_raise TypeError do
      T.check(String, nil)
    end
  end

  def test_types_mismatch
    assert_raise TypeError do
      T.check(String, 123)
    end
  end

  def test_string
    assert_equal "the string", T.check(String, "the string"), "Plain String type should allow String value"
  end

  def test_string_nil
    err = assert_raise TypeError do
      T.check(String, nil)
    end
    assert_match "Type String does not allow nil value", err.message
  end

  def test_boolean
    assert_equal true, T.check(Boolean, true), "Artificial Boolean type should allow true value"
  end

  def test_date
    assert_equal Date.new(2020, 5, 24), T.check(Date, Date.new(2020, 5, 24)), "Date type should pass validation"
  end

  def test_datetime
    assert_equal DateTime.new(2020, 5, 24, 14, 30, 30), T.check(Date, DateTime.new(2020, 5, 24, 14, 30, 30)), "DateTime type should pass validation"
  end

  def test_time
    assert_equal Time.new(2007,11,5,13,45,0, "-05:00"), T.check(Time, Time.new(2007, 11, 5, 13, 45, 0, "-05:00")), "Time type should pass validation"
  end

  def test_uuid
    assert_equal "123e4567-e89b-12d3-a456-426655440000", T.check(UUID, "123e4567-e89b-12d3-a456-426655440000"), "UUID type should pass validation on correctly formatted string"
  end

  def test_uuid_fail
    assert_raise TypeError do
      T.check(UUID, "really not the uuid")
    end
  end

  def test_untyped_success
    assert_equal "bla", T.check(Untyped, "bla"), "Untyped should accept strings"
  end

  def test_untyped_nil
    assert_raise TypeError do
      T.check(Untyped, nil)
    end
  end

  def test_nilable_nil
    assert_equal nil, T.check(T.nilable(String), nil), "Nilable type should allow nil value"
  end

  def test_nilable
    assert_equal "the string", T.check(T.nilable(String), "the string"), "Nilable String type should allow String value"
  end
end

class TypeCheckArray < Test::Unit::TestCase
  def test_array_string
    assert_equal ["the string"], T.check(T.array(String), ["the string"]), "Array of String should allow String value"
  end

  def test_array_fail
    assert_raise TypeError do
      T.check(T.array(String), "the string")
    end
  end

  def test_array_wrong_item_type
    assert_raise TypeError do
      T.check(T.array(String), ["the string", 123])
    end
  end
end

class TypeCheckHash < Test::Unit::TestCase
  def test_hash_string_to_string
    assert_equal({"key" => "the value"}, T.check(T.hash(String, String), {"key" => "the value"}), "Hash of String -> String should allow String -> String value")
  end

  def test_hash_string_to_untyped
    assert_equal({"key" => "the value"}, T.check(T.hash(String, Untyped), {"key" => "the value"}), "Hash of String -> Untyped should allow String -> String value")
  end
end

class TypeCheckAny < Test::Unit::TestCase
  def test_success
    assert_equal(123, T.check(T.any(String, Integer), 123), "Any of String, Integer should allow Integer value")
  end

  def test_fail
    assert_raise TypeError do
      T.check(T.any(String, Integer), true)
    end
  end
end