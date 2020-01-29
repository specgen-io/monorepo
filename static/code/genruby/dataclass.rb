class DataClass
  def self.json_attributes
    @json_attributes
  end

  def self.val(name, type)
    if @json_attributes == nil
      @json_attributes = {}
    end
    @json_attributes[name] = type
    attr_reader name
  end

  def self.var(name, type)
    if @json_attributes == nil
      @json_attributes = {}
    end
    @json_attributes[name] = type
    attr_accessor name
  end

  def self.jsoner_deserialize(json_value)
    T.check(T.hash(String, NilableUntyped), json_value)
    parameters = @json_attributes.map do |attr, attr_type|
      attr_value = json_value[attr.to_s]
      [attr, Jsoner.deserialize(attr_type, attr_value)]
    end
    return self.new parameters.to_h
  end

  def self.jsoner_serialize(value)
    T.check(self, value)
    attrs = @json_attributes.map do |attr, attr_type|
      [attr, Jsoner.serialize(attr_type, value.send(attr))]
    end
    return attrs.to_h
  end

  def initialize(params)
    self.class.json_attributes.each do |attr, attr_type|
      attr_value = params[attr]
      self.instance_variable_set("@#{attr}", T.check(attr_type, attr_value))
    end
  end

  def ==(other)
    begin
      T.check(self.class, other)
      self.class.json_attributes.keys.each do |attr|
        if self.instance_variable_get("@#{attr}") != other.instance_variable_get("@#{attr}")
          return false
        end
      end
      return true
    rescue
      return false
    end
  end
end