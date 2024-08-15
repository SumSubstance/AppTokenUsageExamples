require_relative 'fixed_info'
require_relative 'jsonable'

class Applicant < JSONable
  attr_accessor :fixed_info, :external_user_id, :lang

  # @param [String] external_user_id  - uniq user id in your system
  # @param [Object] lang              - language of commentary for rejection reasons
  def initialize(external_user_id, lang)
    @lang = lang
    @external_user_id = external_user_id
    @fixed_info = FixedInfo.new()
  end

end
