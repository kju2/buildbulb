require 'build-bulb/light'
require 'build-bulb/project'
require 'build-bulb/server'

class Ignore < BasicObject

    def method_missing(name, *args, &block)
        # do nothing
    end

end
