package genruby

import (
	"specgen/gen"
	"strings"
)

func generateBaseClient(moduleName string, path string) *gen.TextFile {
	code := `
require "net/http"
require "net/https"
require "uri"
require "cgi"

require "emery"

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

  module Stringify
    def Stringify.to_string(value)
      if T.instance_of?(DateTime, value)
        value.strftime('%Y-%m-%dT%H:%M:%S')
      else
        value.to_s
      end
    end
  end

  class StringParams
    attr_reader :params

    def initialize
      @params = {}
    end

    def []= (param_name, value)
      if value != nil
        @params[param_name] = Stringify::to_string(value)
      end
    end

    def set(param_name, typedef, value)
      self[param_name] = T.check_var(param_name, typedef, value)
    end

    def query_str
      parts = (@params || {}).map { |param_name, value| "%s=%s" % [param_name, CGI.escape(value)] }
      parts.empty? ? "" : "?"+parts.join("&")
    end

    def set_to_url(url)
      @params.each do |param_name, value|
        url = url.gsub("{#{param_name}}", CGI.escape(value))
      end
      url
    end
  end
end
`
	code, _ = gen.ExecuteTemplate(code, struct{ ModuleName string }{moduleName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}

func generateClientRoot(gemName string, path string) *gen.TextFile {
	code := `
require "emery"
require "[[.GemName]]/models"
require "[[.GemName]]/baseclient"
require "[[.GemName]]/client"`
	code, _ = gen.ExecuteTemplate(code, struct{ GemName string }{gemName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
