lint:
	golangci-lint run
test:
	./scripts/test.sh

mock:
	rm -rf ./mocks && mockery
.PHONY: mock