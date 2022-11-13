require "test/unit/runner/junitxml"

require "test_service_client"

module TestService
  class StringParamsTests < Test::Unit::TestCase
    def test_query_params_set
      params = StringParams.new
      params.set('some_param', Integer, 123)
      params.set('another_param', String, 'bla-bla')
      actual = params.query_str()
      assert_equal "?some_param=123&another_param=bla-bla", actual
    end
  end
end