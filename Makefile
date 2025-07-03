# Copyright 2025 The Nuclio Authors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: fmt
fmt: ensure-golangci-linter
	gofmt -s -w .
	$(GOPATH)/bin/golangci-lint run --fix

.PHONY: lint
lint: modules ensure-golangci-linter
	@echo Linting...
	$(GOPATH)/bin/golangci-lint run -v
	@echo Done.

.PHONY: test
test:
	@echo Running tests...
	go test -v -tags=test_unit -race -coverprofile=coverage.out $(go list ./... | grep '^github.com/nuclio/')
	@echo Done.

.PHONY: test-coverage
test-coverage: test
	@echo Generating coverage report...
	go tool cover -html=coverage.out -o coverage.html
	@echo Coverage report generated: coverage.html

.PHONY: modules
modules: ensure-gopath
	@go mod download

.PHONY: ensure-gopath
ensure-gopath:
ifndef GOPATH
	$(error GOPATH must be set)
endif

.PHONY: clean
clean:
	@echo Cleaning...
	rm -f coverage.out coverage.html
	go clean ./...
	@echo Done.

GOLANGCI_LINT_VERSION := 2.2.0
GOLANGCI_LINT_BIN := $(GOPATH)/bin/golangci-lint
GOLANGCI_LINT_INSTALL_COMMAND := GOBIN=$(GOPATH)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)

.PHONY: ensure-golangci-linter
ensure-golangci-linter:
	@if ! command -v $(GOLANGCI_LINT_BIN) >/dev/null 2>&1; then \
		echo "golangci-lint not found. Installing..."; \
		$(GOLANGCI_LINT_INSTALL_COMMAND); \
	else \
		installed_version=$$($(GOLANGCI_LINT_BIN) version | awk '/version/ {print $$4}' | sed 's/^v//'); \
		if [ "$$installed_version" != "$(GOLANGCI_LINT_VERSION)" ]; then \
			echo "golangci-lint version mismatch ($$installed_version != $(GOLANGCI_LINT_VERSION)). Reinstalling..."; \
			$(GOLANGCI_LINT_INSTALL_COMMAND); \
		fi \
	fi
