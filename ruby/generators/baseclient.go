package generators

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
)

func generateBaseClient(moduleName string, path string) *generator.CodeFile {
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
      elsif T.instance_of?(Date, value)
        value.strftime('%Y-%m-%d')
      else
        value.to_s
      end
    end

    def Stringify.set_params_to_url(url, parameters)
      parameters.each do |param_name, value|
        url = url.gsub("{#{param_name}}", CGI.escape(Stringify::to_string(value)))
      end
      url
    end
  end

  class StringParams
    attr_reader :params

    def initialize
      @params = []
    end

    def set(param_name, typedef, value)
      value = T.check_var(param_name, typedef, value)
      if value != nil
        if value.is_a? Array
          value.each { |item_value| @params.push([param_name, Stringify::to_string(item_value)]) }
        else
          @params.push([param_name, Stringify::to_string(value)])
        end
      end
    end

    def query_str
      parts = (@params || []).map { |param_name, value| "%s=%s" % [param_name, CGI.escape(value)] }
      parts.empty? ? "" : "?"+parts.join("&")
    end
  end
end
`
	code, _ = generator.ExecuteTemplate(code, struct{ ModuleName string }{moduleName})
	return &generator.CodeFile{path, strings.TrimSpace(code)}
}

func generateClientRoot(gemName string, path string) *generator.CodeFile {
	code := `
require "emery"
require "[[.GemName]]/models"
require "[[.GemName]]/baseclient"
require "[[.GemName]]/client"`
	code, _ = generator.ExecuteTemplate(code, struct{ GemName string }{gemName})
	return &generator.CodeFile{path, strings.TrimSpace(code)}
}
