require 'build-bulb/light'
require 'build-bulb/project'
require 'build-bulb/server'

require 'logger'

LOGGER = Logger.new(STDERR)
LOGGER.sev_threshold = Logger::DEBUG

class Ignore < BasicObject

    def method_missing(name, *args, &block)
        # do nothing
    end

end
