task :default => :build
platforms=[:linux_amd64, :darwin_amd64]

PKG_DIST_DIR = ENV["PKG_DIST_DIR"]
if PKG_DIST_DIR.nil? then
	PKG_DIST_DIR="dist"
end

#main build task, pass a platform and version string
task :build, [:platform, :version] => [:clean] do |t, args|
	if args[:platform].nil? then
		abort("No platform defined, aborting!")
	end

	if args[:version].nil? then
		abort("No version defined, aborting!")
	end

	platform = args[:platform].intern
	version = args[:version].intern

	if platforms.include?(platform) == false then
		abort("Unknown/unsupported platform #{platform}")
	end


	puts "Building gomeboycolor (version = #{version}) for #{platform}..."
	case platform
	when :linux_amd64
		Rake::Task["build_linux"].invoke(platform, version)
	when :darwin_amd64
#		Rake::Task["build_darwin"].invoke
	else
		abort("Something bad happened")
	end
end


task :build_linux, [:platform, :version]  do |t, args|
	puts "Packaging linux"
	sh %{CGO_LDFLAGS="-Wl,-Bstatic -lglfw -lGLEW -Wl,-Bdynamic" #{construct_build_command(args[:platform], args[:version])}}
	sh result
end

task :build_darwin, [:platform, :version]  do |t, args|
	puts "Packaging darwin"
	buildcommand = construct_build_command(args[:platform], args[:version])
	sh build_command
	sh "mkdir target/#{args[:platform]}/bin/executable && mv target/#{args[:platform]}/bin/gomeboycolor target/#{args[:platform]}/bin/executable/"
	sh "cp -a #{PKG_DIST_DIR}/darwin_amd64/* target/darwin_amd64/"
end


task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./target/*"
end


def construct_build_command(platform, version) 
	return "go build -o target/#{platform}/bin/gomeboycolor -ldflags=\"-X main.VERSION #{version}\" src/gbc.go src/debugger.go src/config.go"  
end
