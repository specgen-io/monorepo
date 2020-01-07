class Boolean
end

class Any
end

module Type
  class TypeMismatchException < StandardError
    attr_reader :message
    def initialize(message)
      super(message)
      @message = message
    end
  end

  def Type.check_not_nil(value)
    if value == nil
      raise TypeMismatchException.new("Value is nil - class[#{value}] does not allow nil value")
    end
  end

  class TypeDef
  end

  class PlainTypeDef < TypeDef
    attr_reader :value_typedef
    def initialize(plain_type)
      @plain_type = plain_type
    end
    def check(value)
      Type.check_not_nil(value)
      if @plain_type != Any
        if @plain_type == Boolean
          if !value.is_a?(TrueClass) and !value.is_a?(FalseClass)
            raise TypeMismatchException.new("Value type[#{value.class}] - class[TrueClass or FalseClass] is required. value[#{value.inspect.to_s}]")
          end
        elsif @plain_type.included_modules.include?(Ruby::Enum)
          if !@plain_type.value?(value)
            raise TypeMismatchException.new("Value is not member of enum #{@plain_type}. value[#{value.inspect.to_s}]")
          end
        elsif
          if !value.is_a?(@plain_type)
            raise TypeMismatchException.new("Value type[#{value.class}] - class[#{@plain_type}] is required. value[#{value.inspect.to_s}]")
          end
        end
      end
      value
    end
    def to_s
      @plain_type.to_s
    end
  end

  class EnumTypeDef < TypeDef
    attr_reader :enum_class
    def initialize(enum_class)
      @enum_class = enum_class
    end
    def check(value)
      if value != nil
        @value_typedef.check(value)
      end
      value
    end
    def to_s
      enum_class.to_s
    end
  end

  class NillableTypeDef < TypeDef
    attr_reader :value_typedef
    def initialize(value_typedef)
      @value_typedef = value_typedef
    end
    def check(value)
      if value != nil
        @value_typedef.check(value)
      end
      value
    end
    def to_s
      "Nillable[#{@value_typedef.to_s}]"
    end
  end

  class ArrayTypeDef < TypeDef
    attr_reader :item_typedef
    def initialize(item_typedef)
      @item_typedef = item_typedef
    end
    def check(values)
      Type.check_not_nil(values)
      if !values.is_a?(Array)
        raise TypeMismatchException.new("Value type[#{values.class}] - class[Array] is required. value[#{values.inspect.to_s}]")
      end
      values.each { |value| @item_typedef.check(value) }
      values
    end
    def to_s
      "Array[#{@item_typedef.to_s}]"
    end
  end

  class HashTypeDef < TypeDef
    attr_reader :key_typedef
    attr_reader :value_typedef
    def initialize(key_typedef, value_typedef)
      @key_typedef = key_typedef
      @value_typedef = value_typedef
    end
    def check(values)
      Type.check_not_nil(values)
      if !values.is_a?(Hash)
        raise TypeMismatchException.new("Value type[#{values.class}] - class[Hash] is required. value[#{values.inspect.to_s}]")
      end
      values.each do |key, value|
        @key_typedef.check(key)
        @value_typedef.check(value)
      end
      values
    end
    def to_s
      "Hash[#{@key_typedef.to_s}, #{@value_typedef.to_s}]"
    end

    class ClassTypeDef < TypeDef
      attr_reader :klass
      def initialize(klass)
        @klass = klass
      end
      def check(value)
        if value != @klass
          raise TypeMismatchException.new("Expected type[#{values.class}] - class[Hash] is required. value[#{values.inspect.to_s}]")
        end
      end
      def to_s
        "Class[#{@klass.to_s}]"
      end
    end

  end

  def Type.wrap_plain(typedef)
    if typedef.is_a?(TypeDef)
      return typedef
    else
      return Type.plain(typedef)
    end
  end

  def Type.plain(the_type)
    PlainTypeDef.new(the_type)
  end

  def Type.nillable(value_type)
    NillableTypeDef.new(wrap_plain(value_type))
  end

  def Type.array(item_type)
    ArrayTypeDef.new(wrap_plain(item_type))
  end

  def Type.hash(key_type, value_type)
    HashTypeDef.new(wrap_plain(key_type), wrap_plain(value_type))
  end

  def Type.check(typedef, value)
    begin
      Type.wrap_plain(typedef).check(value)
    rescue TypeMismatchException => e
      raise TypeMismatchException.new("Type check failed, expected type: #{typedef.to_s}, value: #{value}")
    end
  end

  def Type.check_field(field_name, typedef, value)
    begin
      Type.wrap_plain(typedef).check(value)
    rescue TypeMismatchException => e
      raise TypeMismatchException.new("Field #{field_name} type check failed, expected type: #{typedef.to_s}, value: #{value}")
    end
  end

end

module Ruby
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
      # Define an enumerated value.
      #
      # === Parameters
      # [key] Enumerator key.
      # [value] Enumerator value.
      def define(key, value)
        @_enum_hash ||= {}
        @_enums_by_value ||= {}

        validate_key!(key)
        validate_value!(value)

        store_new_instance(key, value)

        if upper?(key.to_s)
          const_set key, value
        else
          define_singleton_method(key) { value }
        end
      end

      def store_new_instance(key, value)
        new_instance = new(key, value)
        @_enum_hash[key] = new_instance
        @_enums_by_value[value] = new_instance
      end

      def const_missing(key)
        raise Ruby::Enum::Errors::UninitializedConstantError, name: name, key: key
      end

      # Iterate over all enumerated values.
      # Required for Enumerable mixin
      def each(&block)
        @_enum_hash.each(&block)
      end

      # Attempt to parse an enum key and return the
      # corresponding value.
      #
      # === Parameters
      # [k] The key string to parse.
      #
      # Returns the corresponding value or nil.
      def parse(k)
        k = k.to_s.upcase
        each do |key, enum|
          return enum.value if key.to_s.upcase == k
        end
        nil
      end

      # Whether the specified key exists in this enum.
      #
      # === Parameters
      # [k] The string key to check.
      #
      # Returns true if the key exists, false otherwise.
      def key?(k)
        @_enum_hash.key?(k)
      end

      # Gets the string value for the specified key.
      #
      # === Parameters
      # [k] The key symbol to get the value for.
      #
      # Returns the corresponding enum instance or nil.
      def value(k)
        enum = @_enum_hash[k]
        enum.value if enum
      end

      # Whether the specified value exists in this enum.
      #
      # === Parameters
      # [k] The string value to check.
      #
      # Returns true if the value exists, false otherwise.
      def value?(v)
        @_enums_by_value.key?(v)
      end

      # Gets the key symbol for the specified value.
      #
      # === Parameters
      # [v] The string value to parse.
      #
      # Returns the corresponding key symbol or nil.
      def key(v)
        enum = @_enums_by_value[v]
        enum.key if enum
      end

      # Returns all enum keys.
      def keys
        @_enum_hash.values.map(&:key)
      end

      # Returns all enum values.
      def values
        @_enum_hash.values.map(&:value)
      end

      def to_h
        Hash[@_enum_hash.map do |key, enum|
          [key, enum.value]
        end]
      end

      private

      def upper?(s)
        !/[[:upper:]]/.match(s).nil?
      end

      def validate_key!(key)
        return unless @_enum_hash.key?(key)

        raise Ruby::Enum::Errors::DuplicateKeyError, name: name, key: key
      end

      def validate_value!(value)
        return unless @_enums_by_value.key?(value)

        raise Ruby::Enum::Errors::DuplicateValueError, name: name, value: value
      end
    end
  end
end