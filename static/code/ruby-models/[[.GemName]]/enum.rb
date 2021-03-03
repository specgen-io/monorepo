module [[.ModuleName]]
  module Enum
    attr_reader :key, :value

    def initialize(key, value)
      @key = key
      @value = value
    end

    def self.included(base)
      base.extend Enumerable
      base.extend ClassMethods

      base.private_class_method(:new)
    end

    module ClassMethods
      def check(value)
        T.check_not_nil(self, value)
        if !value?(value)
          raise TypeError.new("Value '#{value.inspect.to_s}' is not a member of enum #{self}")
        end
      end

      def jsoner_deserialize(json_value)
        T.check(self, json_value)
      end

      def jsoner_serialize(value)
        T.check(self, value)
      end

      def define(key, value)
        @_enum_hash ||= {}
        @_enums_by_value ||= {}

        if @_enum_hash.key?(key) then
          raise TypeError.new("Duplicate key: #{key}")
        end

        if @_enums_by_value.key?(value) then
          raise TypeError.new("Duplicate value: #{value}")
        end

        new_instance = new(key, value)
        @_enum_hash[key] = new_instance
        @_enums_by_value[value] = new_instance

        if key.to_s == key.to_s.upcase
          const_set key, value
        else
          define_singleton_method(key) { value }
        end
      end

      def each(&block)
        @_enum_hash.each(&block)
      end

      def parse(k)
        k = k.to_s.upcase
        each do |key, enum|
          return enum.value if key.to_s.upcase == k
        end
        nil
      end

      def key?(k)
        @_enum_hash.key?(k)
      end

      def value(k)
        enum = @_enum_hash[k]
        enum.value if enum
      end

      def value?(v)
        @_enums_by_value.key?(v)
      end

      def key(v)
        enum = @_enums_by_value[v]
        enum.key if enum
      end

      def keys
        @_enum_hash.values.map(&:key)
      end

      def values
        @_enum_hash.values.map(&:value)
      end

      def to_h
        Hash[@_enum_hash.map do |key, enum|
          [key, enum.value]
        end]
      end
    end
  end
end