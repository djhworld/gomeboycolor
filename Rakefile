require 'rbconfig'
task :default => :build

def detect_platform 
	if RbConfig::CONFIG['host_os'] == "mswin32"
		return "windows_#{RbConfig::CONFIG['host_arch']}".downcase
	else
		require 'sys/uname'
		os=Sys::Uname.sysname
		machine=Sys::Uname.machine
		return "#{os}_#{machine}".downcase
	end
end


@build_platform=detect_platform
@version="localbuild"

PKG_DIST_DIR = ENV["PKG_DIST_DIR"]
if PKG_DIST_DIR.nil? then
	PKG_DIST_DIR="dist"
end

#main build task, pass a platform and version string
task :build, [:version] => [:clean] do |t, args|
	if args[:version].nil? then
		@version="localbuild_" + @build_platform
	else
		@version = args[:version] + "_" + @build_platform
	end

	puts "Building gomeboycolor (version = #{@version}) for #{@build_platform}..."

	case @build_platform.intern
	when :windows_x86_64
		Rake::Task["build_windows"].invoke
	when :linux_x86_64
		Rake::Task["build_linux"].invoke
	when :linux_i386
		Rake::Task["build_linux"].invoke
	when :linux_i686
		Rake::Task["build_linux"].invoke
	when :darwin_x86_64
		Rake::Task["build_darwin"].invoke
	else
		abort("Unsupported platform #{@build_platform}")
	end
end

task :build_linux do
	puts "Packaging for #{@build_platform} (static linked binary)"
#	sh "mkdir target/#{@build_platform}/bin"
	sh %{CGO_LDFLAGS="-Wl,-Bstatic -lGLEW -lglfw -Wl,-Bdynamic" #{construct_build_command(@build_platform, @version)}}
end

task :build_darwin, [:version]  do |t, args|
	puts "Packaging for #{@build_platform} (dymanic linked binary)"
	sh construct_build_command(@build_platform, @version)
	sh "mkdir target/#{@build_platform}/bin/executable && mv target/#{@build_platform}/bin/gomeboycolor target/#{@build_platform}/bin/executable/"
	sh "cp -a #{PKG_DIST_DIR}/#{@build_platform}/* target/#{@build_platform}/"
end

task :build_windows, [:version]  do |t, args|
	puts "Packaging for #{@build_platform} (dymanic linked binary)"
	sh construct_build_command(@build_platform, @version)
	sh "cp -a #{PKG_DIST_DIR}/#{@build_platform}/* target/#{@build_platform}/bin/"
end


task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./target/*"
end

def construct_build_command(platform, version) 
	return "go build -a -o target/#{platform}/bin/gomeboycolor -ldflags=\"-X main.VERSION #{version}\" src/gbc.go src/debugger.go src/config.go"  
end
