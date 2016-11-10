PROJECT_ROOT := github.com/sectioneight/md-to-godoc

PKG_FILES = *.go render

.PHONY: dependencies
dependencies:
	@which overalls || go get -u github.com/go-playground/overalls
	@which goveralls || go get -u -f github.com/mattn/goveralls

LINT_LOG := lint.log

.PHONY: lint
lint:
	gofmt -d -s $(PKG_FILES) 2>&1 | tee -a $(LINT_LOG)
	$(foreach dir,$(PKG_FILES),go tool vet $(VET_RULES) $(dir) 2>&1 | tee -a $(LINT_LOG);)
	@[ ! -s $(LINT_LOG) ]

COV_REPORT := overalls.coverprofile

.PHONY: test
test: $(COV_REPORT)

$(COV_REPORT): $(PKG_FILES) $(ALL_SRC)
	overalls -project=$(PROJECT_ROOT) \
		-- -race -v | \
		grep -v "No Go Test Files"

.PHONY: coveralls
coveralls: $(COV_REPORT)
	@goveralls -coverprofile=$(COV_REPORT)

.PHONY: clean
clean:
	rm -f $(COV_REPORT) $(LINT_LOG)
