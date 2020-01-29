module T
  def T.check_not_nil(type, value)
    if value == nil
      raise TypeError.new("Type #{type.to_s} does not allow nil value")
    end
  end

  class UntypedType
    def to_s
      "Untyped"
    end
    def check(value)
      T.check_not_nil(self, value)
    end
  end

  class Nilable
    attr_reader :type
    def initialize(type)
      @type = type
    end
    def to_s
      "Nilable[#{type.to_s}]"
    end
    def check(value)
      if value != nil
        T.check(type, value)
      end
    end
    def ==(other)
      T.instance_of?(Nilable, other) and self.type == other.type
    end
  end

  class AnyType
    attr_reader :types
    def initialize(*types)
      @types = types
    end
    def to_s
      "Any[#{types.map { |t| t.to_s}.join(', ')}]"
    end
    def check(value)
      types.each do |type|
        begin
          T.check(type, value)
          return
        rescue TypeError
        end
      end
      raise TypeError.new("Value '#{value.inspect.to_s}' type is #{value.class} - any of #{@types.map { |t| t.to_s}.join(', ')} required")
    end
    def ==(other)
      T.instance_of?(AnyType, other) and (self.types - other.types).empty?
    end
  end

  class ArrayType
    attr_reader :item_type
    def initialize(item_type)
      @item_type = item_type
    end
    def to_s
      "Array[#{item_type.to_s}]"
    end
    def check(value)
      T.check_not_nil(self, value)
      if !value.is_a? Array
        raise TypeError.new("Value '#{value.inspect.to_s}' type is #{value.class} - Array is required")
      end
      value.each { |item_value| T.check(item_type, item_value) }
    end
    def ==(other)
      T.instance_of?(ArrayType, other) and self.item_type == other.item_type
    end
  end

  class HashType
    attr_reader :key_type
    attr_reader :value_type
    def initialize(key_type, value_type)
      @key_type = key_type
      @value_type = value_type
    end
    def to_s
      "Hash[#{@key_type.to_s}, #{@value_type.to_s}]"
    end
    def check(value)
      T.check_not_nil(self, value)
      if !value.is_a? Hash
        raise TypeError.new("Value '#{value.inspect.to_s}' type is #{value.class} - Hash is required")
      end
      value.each do |item_key, item_value|
        T.check(@key_type, item_key)
        T.check(@value_type, item_value)
      end
    end
    def ==(other)
      T.instance_of?(HashType, other) and self.key_type == other.key_type and self.value_type == other.value_type
    end
  end

  class StringFormatted
    attr_reader :regex
    def initialize(regex)
      @regex = regex
    end
    def to_s
      "String<#@regex>"
    end
    def check(value)
      T.check_not_nil(self, value)
      if !value.is_a? String
        raise TypeError.new("Value '#{value.inspect.to_s}' type is #{value.class} - String is required for StringFormatted")
      end
      if !@regex.match?(value)
        raise TypeError.new("Value '#{value.inspect.to_s}' is not in required format '#{@regex}'")
      end
    end
  end

  def T.check(type, value)
    if type.methods.include? :check
      type.check(value)
    else
      if type != NilClass
        T.check_not_nil(type, value)
      end
      if !value.is_a? type
        raise TypeError.new("Value '#{value.inspect.to_s}' type is #{value.class} - #{type} is required")
      end
    end
    return value
  end

  def T.instance_of?(type, value)
    begin
      T.check(type, value)
      true
    rescue TypeError
      false
    end
  end

  def T.nilable(value_type)
    Nilable.new(value_type)
  end

  def T.array(item_type)
    ArrayType.new(item_type)
  end

  def T.hash(key_type, value_type)
    HashType.new(key_type, value_type)
  end

  def T.any(*typdefs)
    AnyType.new(*typdefs)
  end

  def T.check_field(field_name, type, value)
    begin
      check(type, value)
      return value
    rescue TypeError => e
      raise TypeError.new("Field #{field_name} type check failed, expected type: #{type.to_s}, value: #{value}")
    end
  end
end

Boolean = T.any(TrueClass, FalseClass)
Untyped = T::UntypedType.new
NilableUntyped = T.nilable(Untyped)
UUID = T::StringFormatted.new(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/)