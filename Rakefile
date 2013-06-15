task :default => :build

task :build => :clean do
	puts "Building gomeboycolor..."
	sh %{ go build -o bin/gomeboycolor src/gbc.go src/debugger.go src/config.go  }
end

task :clean do
	puts "Cleaning target dir"
	sh "rm -rf ./bin/*"
end
