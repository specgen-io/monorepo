require "test/unit/runner/junitxml"

require '[[.GemName]]/dataclass'
require '[[.GemName]]/jsoner'

module [[.ModuleName]]
  module Client
    class DataClassTypeEquality < Test::Unit::TestCase
      def test_equals
        assert_true TheClass == TheClass
      end

      def test_not_equals
        assert_false TheClass == TheClassWithNested
      end
    end

    class DataClassFields < Test::Unit::TestCase
      def test_fields_meta
        assert_equal ({:string => String, :int => Integer}), TheClass.json_attributes, "Attributes with types should be available on data class"
      end

      def test_read
        obj = TheClass.new(string: "the string", int: 123)
        assert_equal "the string", obj.string, "Immutable field should be readable"
        assert_equal 123, obj.int, "Mutable field should be readable"
      end

      def test_write_mutable
        obj = TheClass.new(string: "the string", int: 123)
        obj.int = 124
        assert_equal 124, obj.int, "Mutable field should be writable"
      end

      def test_write_immutable
        obj = TheClass.new(string: "the string", int: 123)
        assert_raise NoMethodError do
          obj.string = "the other string"
        end
      end
    end

    class DataClassEquality < Test::Unit::TestCase
      def test_nil
        assert_not_equal nil, TheClass.new(string: "the string", int: 123), "Object should not be equal to nil"
      end

      def test_same_fields_values
        assert_equal TheClass.new(string: "the string", int: 123), TheClass.new(string: "the string", int: 123), "Objects with same fields should be equal"
      end

      def test_different_fields_values
        assert_not_equal TheClass.new(string: "the string 2", int: 123), TheClass.new(string: "the string", int: 123), "Objects with different fields should not be equal"
      end
    end

    class DataClassDeserialization < Test::Unit::TestCase
      def test_deserialize_object
        data = Jsoner.from_json(TheClass, '{"string": "the string", "int": 123}')
        T.check(TheClass, data)
        assert_equal TheClass.new(string: "the string", int: 123), data, "Should parse data class object"
      end

      def test_deserialize_nested_object
        data = Jsoner.from_json(TheClassWithNested, '{"nested": {"string": "the string", "int": 123}}')
        T.check(TheClassWithNested, data)
        assert_equal TheClassWithNested.new(nested: TheClass.new(string: "the string", int: 123)), data, "Should parse nested data class object"
      end

      def test_deserialize_object_fail
        assert_raise JsonerError do
          Jsoner.from_json(TheClass, '"string"')
        end
      end

    end

    class DataClassSerialization < Test::Unit::TestCase
      def test_serialize_object
        assert_equal '{"string":"the string","int":123}', Jsoner.to_json(TheClass, TheClass.new(string: "the string", int: 123)), "nil should be serializable to JSON"
      end

      def test_serialize_array_of_objects
        assert_equal '[{"string":"the string","int":123},{"string":"the string 2","int":456}]', Jsoner.to_json(T.array(TheClass), [TheClass.new(string: "the string", int: 123), TheClass.new(string: "the string 2", int: 456)]), "Array of objects should be serializable to JSON"
      end
    end

    class DataClassCopy < Test::Unit::TestCase
      def test_copy
        a = TheClass.new(string: "the string", int: 123)
        b = a.copy(string: "the other string")
        assert_equal "the string", a.string
        assert_equal "the other string", b.string
      end

      def test_copy_non_existing_field
        assert_raise TypeError do
          a = TheClass.new(string: "the string", int: 123)
          a.copy(non_existing: "the other string")
        end
      end
    end

    class TheClass
      include DataClass

      val :string, String
      var :int, Integer
    end

    class TheClassWithNested
      include DataClass

      val :nested, TheClass
    end
  end
end