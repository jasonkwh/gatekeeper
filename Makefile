default: lint

lint:
	golangci-lint run -c .golangci.yml

test-clean-cache:
	go clean -testcache
