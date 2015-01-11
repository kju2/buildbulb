require 'lifx'

module BuildBulb

    class Color

        GREEN = LIFX::Color.hsb(120, 1, 0.7)
        ORANGE = LIFX::Color.hsb(30, 1, 1)
        RED = LIFX::Color.hsb(0, 1, 1)
        YELLOW = LIFX::Color.hsb(47, 1, 1)

    end

    class LightLocator

        def initialize(logger, label, lifx)
            @logger = logger
            @label = label
            @lifx = lifx
        end

        def find_light
            light = @lifx.lights.with_label(@label)
            if light.nil?
                @logger.error("#{@label} wasn't found. Ignoring commands for #{@label}.")
                Ignore.new
            else
                light
            end
        end

    end

end
