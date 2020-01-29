require "test/unit"

require_relative './type'
require_relative './enum'
require_relative './jsoner'

class TypeCheckEnum < Test::Unit::TestCase
  def test_success
    assert_equal TheEnum::one, T.check(TheEnum, TheEnum::one), "Plain enum type should pass type check"
  end

  def test_fail
    assert_raise TypeError do
      T.check(TheEnum, "non existing")
    end
  end
end

class EnumDeserialization < Test::Unit::TestCase
  def test_deserialize_enum
    data = Jsoner.from_json(TheEnum, '"two"')
    T.check(TheEnum, data)
    assert_equal TheEnum::two, data, "Enum should be parsable from JSON"
  end

  def test_deserialize_enum_non_existing_item
    assert_raise JsonerError do
      Jsoner.from_json(TheEnum, '"non_existing"')
    end
  end
end

class EnumSerialization < Test::Unit::TestCase
  def test_serialize_enum
    json_str = Jsoner.to_json(TheEnum, TheEnum::two)
    assert_equal '"two"', json_str, "Enum should be serializable to JSON"
  end
end

class TheEnum
  include Enum

  define :one, 'one'
  define :two, 'two'
end