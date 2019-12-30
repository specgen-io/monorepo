require "test/unit"

require_relative './typecheck'

class TypeCheckTest < Test::Unit::TestCase
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

  def test_nillable_nil
    assert_equal nil, Type.check(Type.nillable(String), nil), "Nillable type should allow nil value"
  end

  def test_nillable_straight
    assert_equal "the string", Type.check(Type.nillable(String), "the string"), "Nillable String type should allow String value"
  end

  def test_nillable_plain
    assert_equal "the string", Type.check(Type.nillable(Type.plain(String)), "the string"), "Nillable String type should allow String value"
  end

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

  def test_hash_plain
    assert_equal({"key" => "the value"}, Type.check(Type.hash(String, String), {"key" => "the value"}), "Hash of String -> String should allow String -> String value")
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
