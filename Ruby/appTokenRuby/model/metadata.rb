require_relative 'jsonable'

class Metadata < JSONable
  attr_accessor :id_doc_type, :id_doc_sub_type, :country

  # @param [nil] id_doc_sub_type FRONT_SIDE, BACK_SIDE or nil, if you wish to upload only one side.
  # please note, if you set up the side, you should upload both sides.
  def initialize(id_doc_type, country, id_doc_sub_type = nil)
    @id_doc_type = id_doc_type
    @id_doc_sub_type = id_doc_sub_type
    @country = country
  end
end
