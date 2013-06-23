require 'rbconfig'

EXE_NAME="gomeboycolor"

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

currentTag=`git describe --tags`.chomp
@version = currentTag + "_" + @build_platform

task :run, [:ars] do |t, args|
	puts "Running gomeboycolor for #{@build_platform}..."

	case @build_platform.intern
	when :darwin_x86_64
		Rake::Task["run_darwin"].invoke(args[:ars])
	else
		abort("Unsupported platform #{@build_platform}")
	end
end

#main build task, pass a platform and version string
task :build, [:build_type] => [:clean] do |t, args|
	if args[:build_type] != "ci" then
		@version += ".localbuild"
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

#can't get the standard CGO_LDFLAGS to work for linux
task :build_linux => [:setgopath, :get_go_deps] do
	puts "Building for #{@build_platform} (static linked binary)"
	ENV["CGO_LDFLAGS"] = "-Wl,-Bstatic -lGLEW -lglfw -Wl,-Bdynamic"
	puts "Set CGO_LDFLAGS to #{ENV["CGO_LDFLAGS"]}"
	sh %{#{construct_build_command(@build_platform, @version, EXE_NAME)}}
	package(EXE_NAME, @version, "target")
end

task :build_darwin => [:setgopath, :set_cgo_flags, :get_go_deps] do
	puts "Building for #{@build_platform} (dymanic linked binary)"
	sh construct_build_command(@build_platform, @version, EXE_NAME)
	sh "mkdir target/#{EXE_NAME}-#{@version}/bin && mv target/#{EXE_NAME}-#{@version}/#{EXE_NAME} target/#{EXE_NAME}-#{@version}/bin/"
	sh "cp -a #{Dir.pwd}/dist/#{@build_platform}/pkg/* target/#{EXE_NAME}-#{@version}/"
	package(EXE_NAME, @version, "target")
end

task :build_windows => [:setgopath, :set_cgo_flags, :get_go_deps] do |t, args|
	puts "Building for #{@build_platform} (dymanic linked binary)"
	sh construct_build_command(@build_platform, @version, EXE_NAME+".exe")
	sh "cp -a #{Dir.pwd}/dist/#{@build_platform}/pkg/* target/#{EXE_NAME}-#{@version}/"
	package(EXE_NAME, @version, "target")
end

task :run_darwin, [:prog_args] => [:setgopath, :set_cgo_flags, :get_go_deps_no_download] do |t, args|
	ENV["DYLD_LIBRARY_PATH"] = "#{Dir.pwd}/dist/#{@build_platform}/pkg/lib"
	sh "#{construct_run_command} #{args[:prog_args]}"
end

task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./target/*"
end

task :setgopath do
	puts "Setting GOPATH to #{Dir.pwd}"
	ENV["GOPATH"] = Dir.pwd
end

task :set_cgo_flags do
	ENV["CGO_CFLAGS"] = "-I#{Dir.pwd}/dist/#{@build_platform}/include"
	ENV["CGO_LDFLAGS"] = "-L#{Dir.pwd}/dist/#{@build_platform}/lib"

	puts "Set CGO_CFLAGS to #{ENV["CGO_CFLAGS"]}"
	puts "Set CGO_LDFLAGS to #{ENV["CGO_LDFLAGS"]}"
end

task :get_go_deps do
	get_deps(true)
end

task :get_go_deps_no_download do
	get_deps(false)
end

def package(name, version, artifacts_dir)
	filename="#{name}_#{version}.zip"
	puts 
	puts "Packaging up to #{filename}"
	sh "rm -rf artifacts"
	sh "mkdir artifacts"
	cd artifacts_dir
	sh "zip -r ../artifacts/#{filename} ./*"
end

def get_deps(download)
	flag = " "
	if download then
		flag += "-d"
	end

	puts "Getting GO dependencies"
	sh "go get#{flag} code.google.com/p/freetype-go/freetype/truetype"
	sh "go get#{flag} github.com/go-gl/gl"
	sh "go get#{flag} github.com/go-gl/glfw"
end

def construct_build_command(platform, version, exename) 
	return "go build -a -o target/#{exename}-#{version}/#{exename} -ldflags=\"-X main.VERSION #{version}\" src/gbc.go src/debugger.go src/config.go"  
end

def construct_run_command
	return "go run -ldflags=\"-X main.VERSION #{@version}.runlocal\" src/gbc.go src/debugger.go src/config.go"  
end
