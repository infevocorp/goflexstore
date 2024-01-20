lint:
	golangci-lint run

mock:
	rm -rf ./mocks && mockery
.PHONY: mock