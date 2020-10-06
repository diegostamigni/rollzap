

.PHONY: test-ci
test-ci: ## run tests for ci and codecov
	go test -coverprofile=coverage.txt ./...