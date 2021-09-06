require 'json'

class JSONable
  def serialize
    hash = {}
    bad_hash = '{}'
    instance_variables.each do |field_name|
      prepared_field_name = prepare_field_name field_name
      value = instance_variable_get(field_name)
      if value.is_a? JSONable
        serialized_value = value.serialize
        unless serialized_value == bad_hash
          hash[prepared_field_name] = serialized_value
        end
      end

      unless value.is_a? JSONable
        hash[prepared_field_name] = value
      end
    end
    hash.to_json.gsub('\\', '')
  end

  def prepare_value(value)
    test = value.is_a? JSONable


    nil
  end

  def prepare_field_name(field_name)
    prepared_field_name = field_name.to_s[1..-1]
    camel_case_lower(prepared_field_name)
  end

  def camel_case_lower(field_name)
    field_name.split('_').inject([]){ |buffer,e| buffer.push(buffer.empty? ? e : e.capitalize) }.join
  end

  private :prepare_field_name
  private :prepare_value
  private :camel_case_lower
end
