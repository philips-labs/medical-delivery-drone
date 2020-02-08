.PHONY: run
run: build
	@./medical-delivery-drone

.PHONY: build
build: generate
	@go build .

.PHONY: generate
generate:
	@go generate ./...
