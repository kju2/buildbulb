require 'json'
require 'socket'

module BuildBulb

    class Server
        
        TIME_OUT = 600
        OFFICE_HOURS = 6..17

        def initialize(logger, address, port, projects, light)
            @logger = logger
            @address = address
            @port = port
            @projects = projects
            @light = light
        end

        def listen_indefinitely()
            socket = TCPServer.new(@address, @port)
            loop do
                begin
                    @light.set_color(@projects.status, duration: 1)
                    power_light_only_during_office_hours(@light, Time.now, OFFICE_HOURS)

                    @logger.debug("Listen to socket for messages...")
                    readable, _, _ = IO.select([socket], [socket], nil, TIME_OUT)
                    if readable != nil
                        @logger.debug("Receiving message...")
                        msg = read_message(socket.accept)
                        @logger.debug("Received: #{msg}")

                        msg =  JSON.parse(msg)

                        project = msg["name"]
                        status = msg["build"]["status"]
                        @projects.update(project, status)
                    end
                rescue SignalException
                    return
                rescue Exception => e
                    @logger.warn("An error occurred.")
                    @logger.warn(e.message)
                end
            end
        end

        def power_light_only_during_office_hours(light, time, office_hours)
            is_off_time = time.saturday? || time.sunday? || !office_hours.cover?(time.hour)
            if light.on? && is_off_time
                @logger.info("Office hours are over, turn off the light.")
                light.turn_off!
            end

            if light.off? && !is_off_time
                @logger.info("Time to work, turn on the lights.")
                light.turn_on!
            end
        end

        def read_message(connection)
            msg = ""
            while line = connection.gets
                msg << line
            end
            connection.close
            return msg
        end

    end
end
