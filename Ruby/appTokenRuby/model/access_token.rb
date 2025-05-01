# frozen_string_literal: true
require_relative 'jsonable'

class FixedInfo < JSONable
  attr_accessor :user_id, :level_name

  def initialize; end

  def initialize(user_id, level_name)
    @user_id = user_id
    @level_name = level_name
  end

end
