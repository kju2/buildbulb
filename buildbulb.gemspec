Gem::Specification.new do |spec|
    spec.name = "buildbulb"
    spec.version = "0.1.0"
    spec.date = "2014-07-19"
    spec.summary = "Let a LIFX lamp display the current status of jenkins projects."
    spec.authors = ["Kju2"]
    spec.required_ruby_version = ">= 2.0"
    spec.files = ["lib/buildbulb.rb", "lib/buildbulb/light.rb", "lib/buildbulb/project.rb", "lib/buildbulb/server.rb"]
    spec.executables = ["buildbulb", "client"]

    spec.add_runtime_dependency "lifx", "~> 0.4"
    spec.add_runtime_dependency "thor", "~> 0.19"
    spec.add_development_dependency "rspec", "~> 3.0"
end
