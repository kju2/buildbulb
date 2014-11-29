require 'buildbulb/light'
require 'buildbulb/project'
require 'buildbulb/server'

class Ignore < BasicObject

    def method_missing(name, *args, &block)
        # do nothing
    end

end
