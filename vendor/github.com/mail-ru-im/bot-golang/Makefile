all: generate test

$(GOPATH)/bin/easyjson:
	go build -o $(GOPATH)/bin/easyjson github.com/mailru/easyjson/easyjson

$(GOPATH)/bin/golangci-lint:
	go build -o $(GOPATH)/bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: test
test:
	go test -v --cover -coverprofile=cover.out ./...

.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	golangci-lint run

.PHONY: generate
generate: $(GOPATH)/bin/easyjson
	go generate
