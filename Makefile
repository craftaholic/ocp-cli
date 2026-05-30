.PHONY: build install test clean

BINARY_NAME=ocp
INSTALL_PATH=$(HOME)/.local/bin

build:
	go build -o $(BINARY_NAME) .

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installed $(BINARY_NAME) to $(INSTALL_PATH)"
	@echo ""
	@echo "Add shell integration to your shell RC file:"
	@echo ""
	@echo "For zsh (~/.zshrc):"
	@echo '  eval "$$(ocp init hook zsh)"'
	@echo ""
	@echo "For bash (~/.bashrc):"
	@echo '  eval "$$(ocp init hook bash)"'
	@echo ""
	@echo "For fish (~/.config/fish/config.fish):"
	@echo '  ocp init hook fish | source'

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
	go clean

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet
