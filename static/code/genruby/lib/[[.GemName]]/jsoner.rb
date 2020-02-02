require "json"
require "date"

module [[.ModuleName]]
  class JsonerError < StandardError
  end

  module Jsoner
    T::StringFormatted.class_eval do
      def jsoner_deserialize(json_value)
        T.check(self, json_value)
      end
      def jsoner_serialize(value)
        T.check(self, value)
      end
    end

    T::ArrayType.class_eval do
      def jsoner_deserialize(json_value)
        T.check_not_nil(self, json_value)
        if !json_value.is_a?(Array)
          raise JsonerError.new("JSON value type #{json_value.class} is not Array")
        end
        json_value.map { |item_json_value| Jsoner.deserialize(self.item_type, item_json_value) }
      end
      def jsoner_serialize(value)
        if !value.is_a?(Array)
          raise JsonerError.new("Value type #{json_value.class} is not Array")
        end
        value.map { |item| Jsoner.serialize(self.item_type, item) }
      end
    end

    T::HashType.class_eval do
      def jsoner_deserialize(json_value)
        T.check_not_nil(self, json_value)
        if self.key_type != String
          raise JsonerError.new("Hash key type #{self.key_type} is not supported for JSON (de)serialization - key should be String")
        end
        if !json_value.is_a?(Hash)
          raise JsonerError.new("JSON value type #{json_value.class} is not Hash")
        end
        json_value.map do |key, value|
          [T.check(self.key_type, key), Jsoner.deserialize(self.value_type, value)]
        end.to_h
      end
      def jsoner_serialize(value)
        if self.key_type != String
          raise JsonerError.new("Hash key type #{self.key_type} is not supported for JSON (de)serialization - key should be String")
        end
        if !value.is_a?(Hash)
          raise JsonerError.new("Value type #{value.class} is not Hash")
        end
        value.map do |key, value|
          [T.check(self.key_type, key), Jsoner.serialize(self.value_type, value)]
        end.to_h
      end
    end

    T::AnyType.class_eval do
      def jsoner_deserialize(json_value)
        types.each do |type|
          begin
            return Jsoner.deserialize(type, json_value)
          rescue TypeError
          end
        end
        raise JsonerError.new("Value '#{json_value.inspect.to_s}' can not be deserialized as any of #{@types.map { |t| t.to_s}.join(', ')}")
      end
      def jsoner_serialize(value)
        T.check(self, value)
        type = types.find {|t| T.instance_of?(t, value) }
        Jsoner.serialize(type, value)
      end
    end

    T::Nilable.class_eval do
      def jsoner_deserialize(json_value)
        if json_value != nil
          Jsoner.deserialize(self.type, json_value)
        else
          nil
        end
      end
      def jsoner_serialize(value)
        if value != nil
          Jsoner.serialize(self.type, value)
        else
          nil
        end
      end
    end

    module FloatSerializer
      def self.jsoner_deserialize(json_value)
        T.check(T.any(Float, Integer), json_value)
        json_value.to_f
      end
      def self.jsoner_serialize(value)
        T.check(Float, value)
      end
    end

    module DateTimeSerializer
      def self.jsoner_deserialize(json_value)
        T.check(String, json_value)
        begin
          DateTime.strptime(json_value, '%Y-%m-%dT%H:%M:%S')
        rescue
          raise JsonerError.new("Failed to parse DateTime from '#{json_value.inspect.to_s}' format %Y-%m-%dT%H:%M:%S is required")
        end
      end
      def self.jsoner_serialize(value)
        T.check(DateTime, value)
        value.strftime('%Y-%m-%dT%H:%M:%S')
      end
    end

    module DateSerializer
      def self.jsoner_deserialize(json_value)
        T.check(String, json_value)
        begin
          Date.strptime(json_value, '%Y-%m-%d')
        rescue
          raise JsonerError.new("Failed to parse Date from '#{json_value.inspect.to_s}' format %Y-%m-%d is required")
        end
      end
      def self.jsoner_serialize(value)
        T.check(Date, value)
        value.strftime('%Y-%m-%d')
      end
    end

    @@serializers = {
        Float => FloatSerializer,
        Date => DateSerializer,
        DateTime => DateTimeSerializer
    }
    def self.add_serializer(type, serializer)
      @@serializers[type] = serializer
    end

    def Jsoner.from_json(type, json)
      data = JSON.parse(json)
      return deserialize(type, data)
    end

    def Jsoner.deserialize(type, json_value)
      begin
        if type.methods.include? :jsoner_deserialize
          return type.jsoner_deserialize(json_value)
        elsif @@serializers.include? type
          return @@serializers[type].jsoner_deserialize(json_value)
        else
          if ![String, Float, Integer, TrueClass, FalseClass, NilClass].include? type
            raise JsonerError.new("Type #{type} is not supported in Jsoner deserialization")
          end
          return T.check(type, json_value)
        end
      rescue StandardError => error
        raise JsonerError.new(error.message)
      end
    end

    def Jsoner.to_json(type, value)
      JSON.dump(serialize(type, value))
    end

    def Jsoner.serialize(type, value)
      begin
        if type.methods.include? :jsoner_serialize
          return type.jsoner_serialize(value)
        elsif @@serializers.include? type
          return @@serializers[type].jsoner_serialize(value)
        else
          if ![String, Float, Integer, TrueClass, FalseClass, NilClass].include? type
            raise JsonerError.new("Type #{type} is not supported in Jsoner serialization")
          end
          return T.check(type, value)
        end
      rescue StandardError => error
        raise JsonerError.new(error.message)
      end
    end
  end
end