all: build

.PHONY: build
build:
	go build ./cmd/vulners

.PHONY: test
test:
	go test -v ./...

.PHONY: openapi-codegen
openapi-codegen:
	go install -v github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen --config=docs/dto-config.yaml docs/openapi.yaml
	oapi-codegen --config=docs/client-config.yaml docs/openapi.yaml
