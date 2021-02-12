# To show commands and test output, run `make Q= <target>` (empty Q)
Q=@

export GO111MODULE ?= on

.PHONY: deps
deps: go.mod
	$Qgo mod download

.PHONY: vendor
vendor: go.mod
	$Qgo mod vendor

GOTESTFLAGS = -race
ifndef Q
GOTESTFLAGS += -v
endif

.PHONY: test
test:
	$Qgo install $(GOTESTFLAGS)
	$Qgo test $(GOTESTFLAGS) -coverpkg="./..." -coverprofile=.coverprofile $(shell go list ./... | grep -v 'examples')
	# We want to add these back into coverage at some point.
	$Qgrep -vE '(cmd)' < .coverprofile > .covprof && mv .covprof .coverprofile
	$Qgo tool cover -func=.coverprofile

.PHONY: errcheck
errcheck: $(GOPATH)/bin/errcheck
	$Qerrcheck -ignoretests -ignore 'Close' ./...

.PHONY: fmtcheck
fmtcheck: $(GOPATH)/bin/goimports
	$Qecho "run make fmtfix if this fails"
	$Qexit $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.cache/*" | xargs goimports -l | wc -l)

.PHONY: fmtfix
fmtfix: $(GOPATH)/bin/goimports
	$Qgoimports -w $(shell find . -iname '*.go' | grep -v vendor)

.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	$Qgolangci-lint run --skip-files '_test.go$$' --timeout 3m $(find . -type d -not -path '*/\.*' -not -path './vendor/*')

