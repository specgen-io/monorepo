require 'net/http'
require 'net/https'
require 'uri'
require 'cgi'

require 'echo_client/type'

module [[.ModuleName]]
  class BaseClient
    attr_reader :uri
    attr_reader :client

    def initialize(base_uri)
      @base_uri = T.check(URI, base_uri)
      @client = Net::HTTP.new(@base_uri.host, @base_uri.port)
      if @base_uri.scheme == "https"
        @client.use_ssl = true
        @client.verify_mode = OpenSSL::SSL::VERIFY_PEER
        @client.cert_store = OpenSSL::X509::Store.new
        @client.cert_store.set_default_paths
      end
    end
  end

  class StringParams
    attr_reader :params

    def initialize
      @params = {}
    end

    def value_to_s(value)
      if T.instance_of?(DateTime, value)
        value.strftime('%Y-%m-%dT%H:%M:%S')
      else
        value.to_s
      end
    end

    def []= (param_name, value)
      if value != nil
        @params[param_name] = value_to_s(value)
      end
    end

    def query_str
      parts = (@params || {}).map { |k,v| "%s=%s" % [k, CGI.escape(v)] }
      parts.empty? ? "" : "?"+parts.join("&")
    end
  end
end
