module BuildBulb

    class Project

        def initialize(id, expected_status, actual_status, last_updated)
            @id = id
            @expected_status = expected_status ? Status[expected_status] : Status::UNKNOWN
            @actual_status = actual_status ? Status[actual_status] : Status::UNKNOWN
            @last_updated = last_updated ? last_updated : 0
        end

        def id
            return @id
        end

        def actual_status
            return @actual_status
        end

        def actual_status=(status)
            @actual_status = status
            if @last_updated > 0
                @last_updated = Time.now.to_i
            end
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

        def to_json
            return {:id => @id, :expected_status => @expected_status.name, :actual_status => @actual_status.name, :last_updated => @last_updated}
        end

    end

    class Projects

        def initialize(logger, marshaller)
            @logger = logger
            @marshaller = marshaller
            @projects = @marshaller.load
            @marshaller.dump(@projects)
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
                    @logger.info("project \"#{id}\" caused status: \"#{combined_status.name}\".")
                end
            end
            return combined_status
        end

        def update(project, status)
            status = Status[status]

            if !@projects.has_key?(project)
                raise KeyError, "Project \"#{project}\" not found."
            end
            @projects[project].actual_status = status

            @marshaller.dump(@projects)
        end
    end
    
    class ProjectsMemoryMarshaller
        def initialize(projects)
            @projects = projects
        end

        def load
            return @projects
        end

        def dump(projects)
            # do nothing
        end

    end

    class ProjectsFileMarshaller
    
        def initialize(logger, filename)
            @logger = logger
            @filename = filename
        end

        def load
            projects = {}

            if !@filename
                return projects
            end

            json = JSON.parse(File.read(@filename))
            json.each do |project|
                id = project["id"]
                expected_status = project["expected_status"]
                actual_status = project["actual_status"]
                last_updated = project["last_updated"]

                projects[id] = Project.new(id, expected_status, actual_status, last_updated)
            end
            return projects
        end

        def dump(projects)
            objects = []

            projects.each do |id, project|
                objects.push(project.to_json)
            end

            File.open(@filename, 'w') { |file_handler| file_handler.puts objects.to_json }
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
            if !key
                raise KeyError, "Nil is not a valid status key."
            end
            @status.fetch(key.upcase.to_sym)
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
