require "test/unit/runner/junitxml"

require '[[.GemName]]/type'
require '[[.GemName]]/tod'
require '[[.GemName]]/jsoner'


module [[.ModuleName]]
  module Client
    class TypeCheckTimeOfDay < Test::Unit::TestCase
      def test_success
        assert_equal TimeOfDay.parse("10:40:50"), T.check(TimeOfDay, TimeOfDay.parse("10:40:50"))
      end

      def test_fail
        assert_raise TypeError do
          T.check(DateTime, TimeOfDay.parse("10:40:50"))
        end
      end
    end

    class TimeOfDayDeserialization < Test::Unit::TestCase
      def test_deserialize
        data = Jsoner.from_json(TimeOfDay, '"10:40:50"')
        T.check(TimeOfDay, data)
        assert_equal TimeOfDay.parse("10:40:50"), data
      end

      def test_deserialize_fail
        assert_raise JsonerError do
          Jsoner.from_json(TimeOfDay, '"abc"')
        end
      end
    end

    class TimeOfDaySerialization < Test::Unit::TestCase
      def test_serialize
        json_str = Jsoner.to_json(TimeOfDay, TimeOfDay.parse("10:40:50"))
        assert_equal '"10:40:50"', json_str
      end
    end
  end
end