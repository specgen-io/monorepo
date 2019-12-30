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

  class AnyTypeDef < TypeDef
    def check(value)
      value
    end
    def to_s
      "Any"
    end
  end

  class PlainTypeDef < TypeDef
    attr_reader :value_typedef
    def initialize(plain_type)
      @plain_type = plain_type
    end
    def check(value)
      Type.check_not_nil(value)
      if @plain_type == Boolean
        if !value.is_a?(TrueClass) and !value.is_a?(FalseClass)
          raise TypeMismatchException.new("Value type[#{value.class}] - class[TrueClass or FalseClass] is required. value[#{value.inspect.to_s}]")
        end
      else
        if !value.is_a?(@plain_type)
          raise TypeMismatchException.new("Value type[#{value.class}] - class[#{@plain_type}] is required. value[#{value.inspect.to_s}]")
        end
      end
      value
    end
    def to_s
      @plain_type.to_s
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

  def Type.any
    AnyTypeDef.new
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
