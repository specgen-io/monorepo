require "test/unit"
require "date"

require_relative './typecheck'
require_relative './enum'

class TypeCheckTest < Test::Unit::TestCase
  def test_nil
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

  def test_enum_success
    assert_equal TheEnum::FIRST, T.check(TheEnum, TheEnum::FIRST), "Plain enum type should pass type check"
  end

  def test_enum_fail
    assert_raise TypeError do
      T.check(TheEnum, "non existing")
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

  def test_hash_string_string
    assert_equal({"key" => "the value"}, T.check(T.hash(String, String), {"key" => "the value"}), "Hash of String -> String should allow String -> String value")
  end

  def test_hash_string_untyped
    assert_equal({"key" => "the value"}, T.check(T.hash(String, Untyped), {"key" => "the value"}), "Hash of String -> Untyped should allow String -> String value")
  end

  def test_any_string_integer
    assert_equal(123, T.check(T.any(String, Integer), 123), "Any of String, Integer should allow Integer value")
  end

  def test_any_string_integer_fail
    assert_raise TypeError do
      T.check(T.any(String, Integer), true)
    end
  end
end

class TheEnum
  include Enum

  define :FIRST, 'first'
  define :SECOND, 'second'
end