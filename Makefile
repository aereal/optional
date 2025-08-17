GO               := go
TEST_PACKAGE     := ./... ./sqltest
COVERAGE_PROFILE := ./coverage.out

.PHONY: coverage
coverage:
	$(GO) test \
		-covermode count \
		-coverpkg=./... \
		-coverprofile $(COVERAGE_PROFILE) \
		-v \
		$(TEST_PACKAGE)
