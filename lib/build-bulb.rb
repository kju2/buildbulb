require 'build-bulb/light'
require 'build-bulb/project'

require 'logger'

LOGGER = Logger.new(STDERR)
LOGGER.sev_threshold = Logger::FATAL
