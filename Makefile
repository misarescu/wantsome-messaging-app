build: 
	@mkdir -p ./bin 
	@go build -o ./bin ./cmd/...

run-server:
	bin/server $(ARGS)

run-client:
	bin/client $(ARGS)

clean:
	@rm -rf bin/

rebuild: clean build