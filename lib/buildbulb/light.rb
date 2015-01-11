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
                LightProxy.new(@logger, @label, Ignore.new)
            else
                LightProxy.new(@logger, @label, light)
            end
        end

    end

    class LightProxy

        def initialize(logger, label, light)
            @logger = logger
            @label = label
            @light = light
        end

        # Turns the light on synchronously.
        # @return Light
        def turn_on!
            @logger.debug("#{@label}: turn on!")
            @light.turn_on!
            self
        end

        # Turns the light off synchronously.
        # @return Light
        def turn_off!
            @logger.debug("#{@label}: turn off!")
            @light.turn_off!
            self
        end

        # @return [Boolean] Returns true if device is on.
        def on?
            @light.on?
        end

        # @return [Boolean] Returns true if device is off.
        def off?
            @light.off?
        end

        # Attempts to set the color of the light to `color` asynchronously.
        # @param color [Color] The color to be set.
        # @return Light
        def set_color(color)
            @logger.debug("#{@label} color is set to #{color}.")
            @light.set_color(color)
            self
        end

    end

end
