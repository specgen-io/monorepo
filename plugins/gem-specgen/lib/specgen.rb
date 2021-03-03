require 'os'

def get_os
  if OS.mac? then
    return "darwin"
  end
  if OS.posix? then
    return "linux"
  end
  if OS.windows? then
    return "windows"
  end
end

def get_ext
  if OS.windows? then
    return ".exe"
  else
    return ""
  end
end

def get_specgen_path
  File.join(File.dirname(File.expand_path(__FILE__)), "dist", "#{get_os}_amd64", "specgen#{get_ext}")
end

def specgen_client_ruby(spec_file: "./spec.yaml", generate_path: ".")
  specgen_path = get_specgen_path
  command = "#{specgen_path} client-ruby --spec-file #{spec_file} --generate-path #{generate_path}"
  puts "Executing specgen:"
  sh command
end

def specgen_models_ruby(spec_file: "./spec.yaml", generate_path: ".")
  specgen_path = get_specgen_path
  command = "#{specgen_path} models-ruby --spec-file #{spec_file} --generate-path #{generate_path}"
  puts "Executing specgen:"
  sh command
end