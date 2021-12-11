.PHONY: install-tools install-npm-tools install-go-tools
install-tools: | install-go-tools install-npm-tools

install-npm-tools:
	npm install -g @commitlint/cli @commitlint/config-conventional
	npm install -g standard-version

install-go-tools:
	go mod download -x
	cat dev/tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

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

.PHONY: generate-code generate-go-code
generate-code: | generate-go-code format

generate-go-code:
	go generate ./...

.PHONY: release
release:
	standard-version
