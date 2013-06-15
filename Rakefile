task :default => :build
platforms=[:linux_amd64, :darwin_amd64]
task :build, [:platform] => [:clean] do |t, args|
	if args[:platform].nil? then
		abort("No platform defined, aborting!")
	end
	platform = args[:platform].intern

	if platforms.include?(platform) == false then
		abort("Unknown/unsupported platform #{platform}")
	end

	puts "Building gomeboycolor for #{args[:platform]}..."
	sh %{ go build -o bin/gomeboycolor src/gbc.go src/debugger.go src/config.go  }
end

task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./bin/*"
end
