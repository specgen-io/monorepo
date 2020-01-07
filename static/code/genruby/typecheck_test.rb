require "test/unit"
require "date"

require_relative './typecheck'

class TypePlainTest < Test::Unit::TestCase
  def test_plain_success
    assert_equal "the string", Type.check(String, "the string"), "Plain String type should allow String value"
  end

  def test_plain_nil
    assert_raise Type::TypeMismatchException do
      Type.check(String, nil)
    end
  end

  def test_plain_types_mismatch
    assert_raise Type::TypeMismatchException do
      Type.check(String, 123)
    end
  end

  def test_plain_boolean
    assert_equal true, Type.check(Boolean, true), "Artificial Boolean type should allow true value"
  end

  def test_plain_date
    assert_equal Date.new(2020, 5, 24), Type.check(Date, Date.new(2020, 5, 24)), "Date type should pass validation"
  end

  def test_plain_datetime
    assert_equal DateTime.new(2020, 5, 24, 14, 30, 30), Type.check(Date, DateTime.new(2020, 5, 24, 14, 30, 30)), "DateTime type should pass validation"
  end

  def test_plain_time
    assert_equal Time.new(2007,11,5,13,45,0, "-05:00"), Type.check(Time, Time.new(2007,11,5,13,45,0, "-05:00")), "Time type should pass validation"
  end

  def test_plain_uuid
    assert_equal "123e4567-e89b-12d3-a456-426655440000", Type.check(UUID, "123e4567-e89b-12d3-a456-426655440000"), "UUID type should pass validation on correctly formatted string"
  end

  def test_plain_uuid_fail
    assert_raise Type::TypeMismatchException do
      Type.check(UUID, "really not the uuid")
    end
  end

  def test_plain_enum_success
    assert_equal Enum::FIRST, Type.check(Enum, Enum::FIRST), "Plain enum type should pass type check"
  end

  def test_plain_enum_fail
    assert_raise Type::TypeMismatchException do
      Type.check(Enum, "non existing")
    end
  end

  def test_plain_any_success
    assert_equal "bla", Type.check(Any, "bla"), "Any type should accept strings"
  end

  def test_plain_any_nil
    assert_raise Type::TypeMismatchException do
      Type.check(Any, nil)
    end
  end
end

class TypeNillableTest < Test::Unit::TestCase
  def test_nillable_nil
    assert_equal nil, Type.check(Type.nillable(String), nil), "Nillable type should allow nil value"
  end

  def test_nillable_straight
    assert_equal "the string", Type.check(Type.nillable(String), "the string"), "Nillable String type should allow String value"
  end

  def test_nillable_plain
    assert_equal "the string", Type.check(Type.nillable(Type.plain(String)), "the string"), "Nillable String type should allow String value"
  end
end

class TypeArrayTest < Test::Unit::TestCase
  def test_array_plain
    assert_equal ["the string"], Type.check(Type.array(String), ["the string"]), "Array of String should allow String value"
  end

  def test_array_noarray
    assert_raise Type::TypeMismatchException do
      Type.check(Type.array(String), "the string")
    end
  end

  def test_array_wrong_item_type
    assert_raise Type::TypeMismatchException do
      Type.check(Type.array(String), ["the string", 123])
    end
  end
end

class TypeHashTest < Test::Unit::TestCase
  def test_hash_success
    assert_equal({"key" => "the value"}, Type.check(Type.hash(String, String), {"key" => "the value"}), "Hash of String -> String should allow String -> String value")
  end

  def test_hash_string_any
    assert_equal({"key" => "the value"}, Type.check(Type.hash(String, Any), {"key" => "the value"}), "Hash of String -> Any should allow String -> String value")
  end
end

class TypeStringTest < Test::Unit::TestCase
  def test_plain_string
    assert_equal "String", Type.plain(String).to_s
  end

  def test_nillable_string
    assert_equal "Nillable[String]", Type.nillable(String).to_s
  end

  def test_array_string
    assert_equal "Array[String]", Type.array(String).to_s
  end

  def test_hash_string
    assert_equal "Hash[String, String]", Type.hash(String, String).to_s
  end
end

class Enum
  include Ruby::Enum

  define :FIRST, 'first'
  define :SECOND, 'second'
end