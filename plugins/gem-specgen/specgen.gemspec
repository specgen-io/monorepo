Gem::Specification.new do |s|
  s.name        = 'specgen'
  s.summary     = 'Specgen gem tool'
  s.version     = ENV['VERSION']
  s.files       = Dir.glob("{lib}/**/*")
  s.authors     = ['Vladimir Sapronov']
  s.add_runtime_dependency 'os'  
end