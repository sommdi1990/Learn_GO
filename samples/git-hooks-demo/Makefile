.PHONY: lint test install-hooks uninstall-hooks

lint:
	@echo "running golangci-lint (if installed)"
	@golangci-lint run ./... || true

test:
	go test ./...

install-hooks:
	git config core.hooksPath .githooks

uninstall-hooks:
	git config --unset core.hooksPath || true


