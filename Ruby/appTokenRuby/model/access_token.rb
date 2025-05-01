# frozen_string_literal: true
require_relative 'jsonable'

class FixedInfo < JSONable
  attr_accessor :userId, :levelName

  def initialize; end

  def initialize(userId, levelName)
    @userId = userId
    @levelName = levelName
  end

end
