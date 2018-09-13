all: prepare 
	go install

prepare:
	go mod tidy


wasm: prepare
	GOARCH=wasm GOOS=js go build -o static/wasm/gbc.wasm .


