PKGS := $(shell go list ./...)

check: test lint vet fmt-check

test:
	go test -v $(PKGS)

lint:
	golint $(PKGS)

vet:
	go vet $(PKGS) 2>&1 | \
		grep -v '^#' | \
		grep -v 'has json tag but is not exported$$' | \
		grep -v 'bad syntax for struct tag' | \
		grep .; \
		EXIT_CODE=$$?; \
		if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi

fmt-check:
	gofmt -l -s *.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi; \
	goimports -l *.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	gofmt -w -s *.go
	goimports -w *.go

