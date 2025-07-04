GO_TEST = go test
COVERAGE_DIR = coverage
COVERAGE_FILE = $(COVERAGE_DIR)/coverage.out
HTML_REPORT = $(COVERAGE_DIR)/coverage.html

.PHONY: test
test: $(COVERAGE_DIR)
	@echo "🚀 Running tests..."
	$(GO_TEST) -coverprofile=$(COVERAGE_FILE) ./...

.PHONY: coverage
coverage: $(COVERAGE_DIR)
	@echo "📊 Generating coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(HTML_REPORT)
	@echo "✅ Coverage report generated at $(HTML_REPORT)"

.PHONY: test-all
test-all: test coverage
	@echo "✔️ All tests completed with coverage report"

$(COVERAGE_DIR):
	@mkdir -p $(COVERAGE_DIR)