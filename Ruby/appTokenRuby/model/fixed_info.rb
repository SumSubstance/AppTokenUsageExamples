# frozen_string_literal: true
require_relative 'jsonable'

class FixedInfo < JSONable
  attr_accessor :first_name, :last_name, :middle_name, :dob

  def initialize; end

  # https://developers.sumsub.com/api-reference/#creating-an-applicant
  # please note we do not recommend usage of provided info functional, cause it can drop your conversion rate.
  # @param [nil] first_name   - name of the applicant
  # @param [nil] last_name    - last name of the applicant
  # @param [nil] dob          - date of birth of the applicant in the yyyy-mm-dd format
  # @return [nil]
  def initialize(first_name, last_name, middle_name, dob)
    @first_name = first_name
    @middle_name = middle_name
    @last_name = last_name
    @dob = dob
  end

end
