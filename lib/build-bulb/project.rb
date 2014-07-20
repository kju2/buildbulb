module BuildBulb

    class Project

        def initialize(id, expected_status, actual_status, last_updated)
            @id = id
            @expected_status = expected_status ? Status[expected_status] : Status::SUCCESS
            @actual_status = actual_status ? Status[actual_status] : Status::UNKNOWN
            @last_updated = last_updated ? last_updated : 0
        end

        def actual_status
            return @actual_status
        end

        def actual_status=(status)
            @actual_status = status
            if @last_updated > 0
                @last_updated = Time.now.to_i
            end
            LOGGER.info("#{@id} now has status #{@actual_status.name}.")
        end

        def status
            if 0 < @last_updated && @last_updated < (Time.now - 2 * 24 * 3600).to_i
                @actual_status = Status::UNKNOWN
            end
            if self.actual_status == @expected_status
                Status::SUCCESS
            else
                @actual_status
            end
        end

    end

    class Projects

        def initialize(projects)
            @projects = {}
            projects.each do |id, status|
                @projects[id] = Project.new(id, status["expected_status"], status["actual_status"], status["last_updated"])
            end
        end

        def [](key)
            if @projects.has_key?(key)
                @projects[key]
            else
                nil
            end
        end

        def status
            combined_status = Status::SUCCESS
            @projects.each do |id, project| 
                if project.status.precedence < combined_status.precedence
                    combined_status = project.status
                    LOGGER.info("project \"#{id}\" caused status: \"#{combined_status.name}\".")
                end
            end
            combined_status
        end

    end
    
    class Status

        def initialize(name, value, color)
            @name = name
            @value = value
            @color = color
        end
        
        def name
            @name
        end

        def precedence
            @value
        end

        def color
            @color
        end

        def Status.[](key)
            @status.fetch(key.to_sym)
        end

        def Status.add_item(key, value, color)
            @status ||= {}
            @status[key] = Status.new(key, value, color)
        end

        def Status.const_missing(key)
            @status.fetch(key)
        end

        Status.add_item(:FAILURE, 0, Color::RED)
        Status.add_item(:UNKNOWN, 1, Color::ORANGE)
        Status.add_item(:UNSTABLE, 2, Color::YELLOW)
        Status.add_item(:SUCCESS, 3, Color::GREEN)

    end
end
