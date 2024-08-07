.PHONY: start

bin/gowq:
	@go build -o bin/gowq cmd/main.go

start: bin/gowq
	@bin/gowq
