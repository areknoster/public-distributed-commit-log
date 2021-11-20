install-go-tools:
	go mod download -x
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

format: format-go format-add-trailing-newline

format-go:
	gofumpt -w .

format-add-trailing-newline:
	git grep -zIl ''  | while IFS= read -rd '' f; do tail -c1 < "$$f" | read -r _ || echo >> "$$f"; done

lint:
	golangci-lint run

generate-code: | generate-go-code format

generate-go-code:
	go generate ./...

.PHONY: install-go-tools format format-go format-add-trailing-newline lint generate-go-code generate-code
