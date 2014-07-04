require 'lifx'
require 'socket'
require 'json'

label = ARGV[0]
label = 'BuildOrb001'

address = '0.0.0.0'
port = 4712

success = "SUCCESS"
unstable = "UNSTABLE"
broken = "BROKEN"
unkown = broken

projects = {"notification-plugin" => {:expected_status => success, :actual_status => unkown},
            "test" => {:expected_status => unstable, :actual_status => unkown}}

#lifx = LIFX::Client.lan
#lifx.discover! do |c|
    #c.lights.with_label(label)
#end


#light = if lifx.tags.include?('BuildOrb001')
  #lights = lifx.lights.with_tag('BuildOrb001')
  #if lights.empty?
    #puts "No lights in the Build Light tag, using the first light found."
    #lifx.lights.first
  #else
    #lights
  #end
#else
  #lifx.lights.first
#end

#if !light
  #puts "No LIFX lights found."
  #exit 1
#end

def office_hours
    morning_light_on = 6
    evening_light_off = 20

    t = Time.now
    return !(t.saturday? || t.sunday? || t.hour < morning_light_on || t.hour > evening_light_off)
end

connection = TCPServer.new(address, port)

def update_light(light, status)
    color = case status
    when success
        LIFX::Color.hsb(120, 1, 1)  
    when unstable
        LIFX::Color.hsb(60, 1, 1)  
    when broken
        LIFX::Color.hsb(0, 1, 1)
    end

    light.set_color(color, duration: 0.2)
end

loop do
    puts "#{Time.now}: Check for office hours..."
    #if light.on? && !office_hours
        #puts 'Office hours are over, turn off the lights...'
        #light.turn_off!
    #end

    #if light.off? && office_hours
        #puts 'Time to work, turn on the lights...'
        #light.turn_on!
    #end

    puts "#{Time.now}: Check socket for new input..."
    readable, _, _ = IO.select([connection], [connection], nil, 60)
    if readable != nil
        puts "#{Time.now}: Parsing message..."
        sock = connection.accept
        buff = ""
        while line = sock.gets
            buff << line
        end
        #puts "#{Time.now}: #{buff}"
        x = JSON.parse(buff)
        project_name = x["name"]
        if projects.has_key?(project_name)
            puts "#{Time.now}: project \"#{project_name}\" status was \"#{projects[project_name][:actual_status]}\" and now it is \"#{x["build"]["status"]}\""
            projects[project_name][:actual_status] = x["build"]["status"]
        end
        sock.close
    end
    
    puts "#{Time.now}: Determine new lamp status..."
    lamp_status = success
    projects.each do |project, status| 
        if status[:expected_status] != status[:actual_status] 
            if status[:actual_status] == broken
                puts "#{Time.now}: project \"#{project}\" caused status: \"#{broken}\""
                lamp_status = broken
            elsif status[:actual_status] == unstable && lamp_status != broken
                puts "#{Time.now}: project \"#{project}\" caused status: \"#{unstable}\""
                lamp_status = unstable
            end
        end
    end

    #update_light(light, lamp_status)
end


#puts "Using light(s): #{light}"
#repo_path = ARGV.first || 'rails/rails'

#repo = Travis::Repository.find(repo_path)
#puts "Watching repository #{repo.slug}"

#def update_light(light, repository)
  #color = case repository.color
  #when 'green'
    #LIFX::Color.hsb(120, 1, 1)  
  #when 'yellow'
    #LIFX::Color.hsb(60, 1, 1)  
  #when 'red'
    #LIFX::Color.hsb(0, 1, 1)
  #end

  #light.turn_on!
  #light.set_color(color, duration: 0.2)
  #puts "#{Time.now}: Build ##{repository.last_build.number} is #{repository.color}."
#end

#update_light(light, repo)

#Travis.listen(repo) do |stream|
  #stream.on('build:started', 'build:finished') do |event|
    #update_light(light, event.repository)
  #end
#end

