require 'json'
require 'logger'
require 'socket'

require './light'
require './project'

module BuildOrb

    DEBUG = false
    PORT = 0
    NAME = 1
    PROJECTS = 2

    OFFICE_HOURS = 6..16

    LOGGER = Logger.new(STDOUT)

    def main()
        # address:port for receiving updates
        address = "0.0.0.0"
        port = ARGV[PORT] ? ARGV[PORT] : 4712
        LOGGER.debug("port is #{port}")

        # name of the light to control
        light_label = ARGV[NAME] ? ARGV[NAME] : "BBEBB"
        light = get_light(light_label, debug: DEBUG)
        LOGGER.debug("light is #{light}")

        projects_file_name = ARGV[PROJECTS] ? ARGV[PROJECTS] : "projects.json"
        projects = Projects.new(JSON.parse(File.read(projects_file_name)))
        LOGGER.debug("projects are #{projects}")

        socket = TCPServer.new(address, port)
        LOGGER.info("Listen to socket for messages...")
        loop do
            power_light_only_during_office_hours(light, Time.now, OFFICE_HOURS)

            readable, _, _ = IO.select([socket], [socket], nil, 600)
            if readable != nil
                LOGGER.debug("Receiving message...")
                msg = read_message(socket.accept)
                LOGGER.debug("Received: #{msg}")

                msg =  JSON.parse(msg)
                project = projects[msg["name"]]
                status = msg["build"]["status"]
                if project && status
                    project.actual_status = Status[status]
                end
            end
            
            status = projects.status
            LOGGER.info("the combined status of all projects is \"#{status.name}\".")
            light.set_color(status.color, duration: 0.2)
        end
    end

    def power_light_only_during_office_hours(light, time, office_hours)
        is_off_time = time.saturday? || time.sunday? || !office_hours.cover?(time.hour)
        if light.on? && is_off_time
            LOGGER.info("Office hours are over, turn off the light.")
            light.turn_off!
        end

        if light.off? && !is_off_time
            LOGGER.info("Time to work, turn on the lights.")
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

include BuildOrb
BuildOrb::main
