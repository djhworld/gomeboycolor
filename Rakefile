task :default => :build
platforms=[:linux_amd64, :darwin_amd64]

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

	puts "Building gomeboycolor v#{version} for #{platform}..."
	sh %{ go build -o target/#{platform}/bin/executable/gomeboycolor -ldflags="-X main.VERSION #{version}" src/gbc.go src/debugger.go src/config.go  }

	case platform
	when :linux_amd64
		Rake::Task["pkg_linux"].invoke
	when :darwin_amd64
		puts "Darwin is unsupported!"
	else
		abort("Something bad happened")
	end
end


task :pkg_linux do
	puts "Packaging linux"
	sh "cp -a dist/linux_amd64/* target/linux_amd64/"
end

task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./target/*"
end
