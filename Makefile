BUILD=go build

default: protoc client

client:
	$(BUILD) -o ./bin/client ./cmd/client
.PHONY: client

protoc:
	./scripts/protoc.sh
.PHONY: protoc
