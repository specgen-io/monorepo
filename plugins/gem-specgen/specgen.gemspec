Gem::Specification.new do |s|
  s.name        = 'specgen'
  s.summary     = 'Spec gem tool'
  s.version     = ENV['VERSION']
  s.files       = Dir.glob("{lib}/**/*")
  s.authors     = ['Moda Operandi']
  s.add_runtime_dependency 'os'  
end