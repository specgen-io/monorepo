require 'net/http'
require 'net/https'
require 'uri'

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
end
