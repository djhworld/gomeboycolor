all: prepare 
	go install

prepare:
	go mod tidy


