all:
	go build -o bin/gomeboycolor src/gbc.go src/debugger.go src/config.go

clean:
	echo "Cleaning bin dir"
	rm ./bin/*
