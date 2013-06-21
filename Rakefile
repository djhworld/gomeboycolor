require 'rbconfig'
task :default => :build

def detect_platform 
	if RbConfig::CONFIG['host_os'] == "mswin32"
		return "windows_#{RbConfig::CONFIG['host_cpu']}".downcase
	else
		require 'sys/uname'
		os=Sys::Uname.sysname
		machine=Sys::Uname.machine
		return "#{os}_#{machine}".downcase
	end
end


@build_platform=detect_platform
@version="localbuild"

EXE_NAME="gomeboycolor"
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

task :build_linux => [:setgopath, :get_go_deps] do
	puts "Packaging for #{@build_platform} (static linked binary)"
	sh %{CGO_LDFLAGS="-Wl,-Bstatic -lGLEW -lglfw -Wl,-Bdynamic" #{construct_build_command(@build_platform, @version, EXE_NAME)}}
end

task :build_darwin do
	puts "Packaging for #{@build_platform} (dymanic linked binary)"
	sh construct_build_command(@build_platform, @version, EXE_NAME)
	sh "mkdir target/#{@build_platform}/bin && mv target/#{@build_platform}/#{EXE_NAME} target/#{@build_platform}/bin/"
	sh "cp -a #{PKG_DIST_DIR}/#{@build_platform}/* target/#{@build_platform}/"
end

task :build_windows => [:setgopath, :get_go_deps] do |t, args|
	puts "Packaging for #{@build_platform} (dymanic linked binary)"
	ENV["CGO_CFLAGS"] = "-I#{Dir.pwd}/dist/#{@build_platform}/include"
	ENV["CGO_LDFLAGS"] = "-L#{Dir.pwd}/dist/#{@build_platform}/lib"
	sh construct_build_command(@build_platform, @version, EXE_NAME+".exe")
	sh "cp -a #{Dir.pwd}/dist/#{@build_platform}/pkg/* target/#{@build_platform}/"
end


task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./target/*"
end

task :setgopath do
	puts "Setting GOPATH to #{Dir.pwd}"
	ENV["GOPATH"] = Dir.pwd
end

task :get_go_deps do
	puts "Getting GO dependencies"
	sh "go get -d code.google.com/p/freetype-go/freetype/truetype"
	sh "go get -d github.com/go-gl/gl"
	sh "go get -d github.com/go-gl/glfw"
end

def construct_build_command(platform, version, exename) 
	return "go build -a -o target/#{platform}/#{exename} -ldflags=\"-X main.VERSION #{version}\" src/gbc.go src/debugger.go src/config.go"  
end
