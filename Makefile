# --------- BUILD & RELEASE ------------
COMMIT = $(shell git rev-parse HEAD)
.PHONY: build build-docker build-linux
build-linux:
	mkdir -p bin/linux
	GOARCH=amd64 GOOS=linux go build -o bin/linux ./cmd/...

build-darwin:
	mkdir -p bin/darwin
	GOARCH=amd64 GOOS=darwin go build -o bin/darwin ./cmd/...

build: build-linux build-darwin build-docker

build-docker:
	docker build -t pdcl-acceptance-sentinel:$(COMMIT) .
	docker tag pdcl-acceptance-sentinel:$(COMMIT) pdcl-acceptance-sentinel:latest

.PHONY: release
release:
	standard-version -s

clean:
	rm bin/*

# --------- TESTS ------------
.PHONY: unit-test
unit-test:
	go test -race -short -count=1 ./...

# --------- TOOLS ------------
.PHONY: install-tools install-npm-tools install-go-tools
install-tools: | install-go-tools install-npm-tools
install-npm-tools:
	npm install -g @commitlint/cli @commitlint/config-conventional
	npm install -g standard-version

install-go-tools:
	go mod download -x
	cat dev/tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

# --------- FORMAT & LINT ------------
.PHONY: format format-go format-add-trailing-newline
format: format-go format-add-trailing-newline

format-go:
	gofumpt -w .
	gci -w -local github.com/areknoster/public-distributed-commit-log . 1>/dev/null

format-add-trailing-newline:
	git grep -zIl ''  | while IFS= read -rd '' f; do tail -c1 < "$$f" | read -r _ || echo >> "$$f"; done

.PHONY: lint lint-go lint-commits
lint: | lint-go lint-commits

lint-go:
	golangci-lint run

lint-commits:
	commitlint --from main --config commitlint.config.yaml


# --------- CODEGEN ------------
.PHONY: generate-code generate-go-code
generate-code: | generate-go-code format

generate-go-code:
	go generate ./...
