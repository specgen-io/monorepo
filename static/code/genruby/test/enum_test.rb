require "test/unit/runner/junitxml"

require '[[.GemName]]/type'
require '[[.GemName]]/enum'
require '[[.GemName]]/jsoner'

module [[.ModuleName]]
  class TypeCheckEnum < Test::Unit::TestCase
    def test_success
      assert_equal SomeEnum::one, T.check(SomeEnum, SomeEnum::one), "Plain enum type should pass type check"
    end

    def test_fail
      assert_raise TypeError do
        T.check(SomeEnum, "non existing")
      end
    end
  end

  class EnumDeserialization < Test::Unit::TestCase
    def test_deserialize_enum
      data = Jsoner.from_json(SomeEnum, '"two"')
      T.check(SomeEnum, data)
      assert_equal SomeEnum::two, data, "Enum should be parsable from JSON"
    end

    def test_deserialize_enum_non_existing_item
      assert_raise JsonerError do
        Jsoner.from_json(SomeEnum, '"non_existing"')
      end
    end
  end

  class EnumSerialization < Test::Unit::TestCase
    def test_serialize_enum
      json_str = Jsoner.to_json(SomeEnum, SomeEnum::two)
      assert_equal '"two"', json_str, "Enum should be serializable to JSON"
    end
  end

  class SomeEnum
    include Enum

    define :one, 'one'
    define :two, 'two'
  end
end