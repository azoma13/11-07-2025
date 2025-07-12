.PHONY:
.SILENT:

build: 
	go build -o ./.bin/archiving-service ./cmd/app/main.go

run: build 
	./.bin/archiving-service