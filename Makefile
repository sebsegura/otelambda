BIN= $(CURDIR)/bin
FUNCTION=create-contact-aws-lambda

.PHONY: build

build: fmt
	@env GOOS=linux go build -ldflags="-s -w" -o bin/create/bootstrap cmd/create/main.go
	@env GOOS=linux go build -ldflags="-s -w" -o bin/update/bootstrap cmd/update/main.go

lint:
	@go vet ./...

fmt:
	@go fmt ./...

clean:
	@rm -rf $(BIN)

zip: build
	@zip -j $(FUNCTION).zip bin/bootstrap

.PHONY: test
test:
	@go test -v -cover ./...