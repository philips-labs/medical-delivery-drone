.PHONY: run
run: build
	@./medical-delivery-drone

.PHONY: build
build:
	@go build .
