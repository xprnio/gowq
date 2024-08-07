.PHONY: start

bin/work-queue:
	@go build -o bin/work-queue cmd/main.go

start: bin/work-queue
	@bin/work-queue
