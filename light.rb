require 'lifx'

module BuildOrb

    class Color

        GREEN = LIFX::Color.hsb(120, 1, 1)
        ORANGE = LIFX::Color.hsb(60, 1, 1) # TODO find orange
        RED = LIFX::Color.hsb(0, 1, 1)
        YELLOW = LIFX::Color.hsb(60, 1, 1) # TODO ask colleagues which one the want

    end

    class FakeLight

        def initialize(id)
            @id = id
            @status = :off
        end

        # Turns the light(s) on synchronously
        # @return [Light, LightCollection] self for chaining
        def turn_on!
            LOGGER.info("#{@id}: turn on!")
            @status = :on
            self
        end

        # Turns the light(s) off synchronously
        # @return [Light, LightCollection]
        def turn_off!
            LOGGER.info("#{@id}: turn off!")
            @status = :off
            self
        end

        # @see #power
        # @return [Boolean] Returns true if device is on
        def on?(refresh: false, fetch: true)
            @status == :on
        end

        # @see #power
        # @return [Boolean] Returns true if device is off
        def off?(refresh: false, fetch: true)
            @status == :off
        end

        # Attempts to set the color of the light(s) to `color` asynchronously.
        # @param color [Color] The color to be set
        # @param duration: [Numeric] Transition time in seconds
        # @return [Light, LightCollection] self for chaining
        def set_color(color, duration: LIFX::Config.default_duration)
            self
        end

    end

    def get_light(id, debug: false)
        if debug
            return BuildOrb::FakeLight.new(id)
        else
            LIFX::Config.logger.level = Logger::DEBUG
            lifx = LIFX::Client.lan
            lifx.discover! do |c|
                c.lights.with_label(id)
            end

            return lifx.lights.with_label(id)
        end
    end

end
