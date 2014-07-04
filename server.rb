require 'lifx'
require 'socket'
require 'json'

label = ARGV[0]
label = 'BBEBB'

address = '0.0.0.0'
port = 4712

SUCCESS = "SUCCESS"
UNSTABLE = "UNSTABLE"
BROKEN = "BROKEN"
UNKOWN = BROKEN

projects = {"notification-plugin" => {:expected_status => SUCCESS, :actual_status => UNKOWN},
            "test" => {:expected_status => UNSTABLE, :actual_status => UNKOWN}}

lifx = LIFX::Client.lan
lifx.discover! do |c|
    c.lights.with_label(label)
end


LIGHT = if lifx.tags.include?(label)
  lights = lifx.lights.with_tag(label)
  if lights.empty?
    puts "No lights in the Build Light tag, using the first light found."
    lifx.lights.first
  else
    lights
  end
else
  lifx.lights.first
end

if !LIGHT
  puts "No LIFX lights found."
  exit 1
end

def office_hours
    morning_light_on = 6
    evening_light_off = 20

    t = Time.now
    return !(t.saturday? || t.sunday? || t.hour < morning_light_on || t.hour > evening_light_off)
end

connection = TCPServer.new(address, port)

def update_light(light, status)
    color = case status
    when SUCCESS
        LIFX::Color.hsb(120, 1, 1)  
    when UNSTABLE
        LIFX::Color.hsb(60, 1, 1)  
    when BROKEN
        LIFX::Color.hsb(0, 1, 1)
    end

    light.set_color(color, duration: 0.2)
end

loop do
    puts "#{Time.now}: Check for office hours..."
    if LIGHT.on? && !office_hours
        puts 'Office hours are over, turn off the lights...'
        LIGHT.turn_off!
    end

    if LIGHT.off? && office_hours
        puts 'Time to work, turn on the lights...'
        LIGHT.turn_on!
    end

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
    lamp_status = SUCCESS
    projects.each do |project, status| 
        if status[:expected_status] != status[:actual_status] 
            if status[:actual_status] == BROKEN
                puts "#{Time.now}: project \"#{project}\" caused status: \"#{BROKEN}\""
                lamp_status = BROKEN
            elsif status[:actual_status] == UNSTABLE && lamp_status != BROKEN
                puts "#{Time.now}: project \"#{project}\" caused status: \"#{UNSTABLE}\""
                lamp_status = UNSTABLE
            end
        end
    end

    update_light(LIGHT, lamp_status)
end
