require "test/unit/runner/junitxml"

require "test_service_client"

module TestService
  class ClientTests < Test::Unit::TestCase
    def test_echo_body_string
      request_body = "some text"
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_body_string(body: request_body)
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal request_body, response.ok
    end

    def test_echo_body_model
      request_body = Message.new(int_field: 123, string_field: "the string")
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_body_model(body: request_body)
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal request_body, response.ok
    end

    def test_echo_body_array
      request_body = ["the str1", "the str2"]
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_body_array(body: request_body)
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal request_body, response.ok
    end

    def test_echo_body_map
      request_body = {"one" => "the str1", "two" => "the str2"}
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_body_map(body: request_body)
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal request_body, response.ok
    end

    def test_echo_query_parameters
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_query(
        int_query: 123,
        long_query: 123,
        float_query: 1.23,
        double_query: 1.23,
        decimal_query: 1.23,
        bool_query: true,
        string_query: "the string",
        string_opt_query: "the optional string",
        string_defaulted_query: "the defaulted string",
        string_array_query: ["the string array item"],
        uuid_query: "123e4567-e89b-12d3-a456-426655440000",
        date_query: Date.new(2019, 11, 30),
        date_array_query: [Date.new(2019, 11, 30)],
        datetime_query: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_query: "SECOND_CHOICE"
      )
      expected = Parameters.new(
        int_field: 123,
        long_field: 123,
        float_field: 1.23,
        double_field: 1.23,
        decimal_field: 1.23,
        bool_field: true,
        string_field: "the string",
        string_opt_field: "the optional string",
        string_defaulted_field: "the defaulted string",
        string_array_field: ["the string array item"],
        uuid_field: "123e4567-e89b-12d3-a456-426655440000",
        date_field: Date.new(2019, 11, 30),
        date_array_field: [Date.new(2019, 11, 30)],
        datetime_field: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_field: "SECOND_CHOICE"
      )
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal expected, response.ok
    end

    def test_echo_header
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_header(
        int_header: 123,
        long_header: 123,
        float_header: 1.23,
        double_header: 1.23,
        decimal_header: 1.23,
        bool_header: true,
        string_header: "the string",
        string_opt_header: "the optional string",
        string_defaulted_header: "the defaulted string",
        string_array_header: ["the string array item"],
        uuid_header: "123e4567-e89b-12d3-a456-426655440000",
        date_header: Date.new(2019, 11, 30),
        date_array_header: [Date.new(2019, 11, 30)],
        datetime_header: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_header: "SECOND_CHOICE"
      )
      expected = Parameters.new(
        int_field: 123,
        long_field: 123,
        float_field: 1.23,
        double_field: 1.23,
        decimal_field: 1.23,
        bool_field: true,
        string_field: "the string",
        string_opt_field: "the optional string",
        string_defaulted_field: "the defaulted string",
        string_array_field: ["the string array item"],
        uuid_field: "123e4567-e89b-12d3-a456-426655440000",
        date_field: Date.new(2019, 11, 30),
        date_array_field: [Date.new(2019, 11, 30)],
        datetime_field: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_field: "SECOND_CHOICE"
      )
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal expected, response.ok
    end

    def test_echo_url_params
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_url_params(
        int_url: 123,
        long_url: 123,
        float_url: 1.23,
        double_url: 1.23,
        decimal_url: 1.23,
        bool_url: true,
        string_url: "the-string",
        uuid_url: "123e4567-e89b-12d3-a456-426655440000",
        date_url: Date.new(2019, 11, 30),
        datetime_url: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_url: "SECOND_CHOICE"
      )
      assert_true response.respond_to? :ok
      assert_true response.ok?
      expected = UrlParameters.new(
        int_field: 123,
        long_field: 123,
        float_field: 1.23,
        double_field: 1.23,
        decimal_field: 1.23,
        bool_field: true,
        string_field: "the-string",
        uuid_field: "123e4567-e89b-12d3-a456-426655440000",
        date_field: Date.new(2019, 11, 30),
        datetime_field: DateTime.new(2019, 11, 30, 17, 45, 55),
        enum_field: "SECOND_CHOICE"
      )
      assert_equal expected, response.ok
    end
  end
end

module TestService::V2
  class ClientTests < Test::Unit::TestCase
    def test_echo_body_model
      request_body = Message.new(bool_field: true, string_field: "the string")
      client = EchoClient.new(URI(ENV["SERVICE_URL"]))
      response = client.echo_body_model(body: request_body)
      assert_true response.respond_to? :ok
      assert_true response.ok?
      assert_equal request_body, response.ok
    end
  end
end