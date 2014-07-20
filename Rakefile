require "rake"
require "rspec/core/rake_task"

RSpec::Core::RakeTask.new(:spec) do |t|
  t.rspec_opts = "--color --format doc"
  t.verbose = false
end

task :default => :spec
