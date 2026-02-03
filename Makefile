.PHONY: build install uninstall test ext-build ext-package ext-install

BINARY_NAME := clari
INSTALL_PATH := /usr/local/bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/claritask
	@echo "Built $(BINARY_NAME)"

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_PATH)/$(BINARY_NAME)"

uninstall:
	@rm -f $(BINARY_NAME)
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

test:
	@go test ./test/... -v

# VSCode Extension
ext-build:
	@echo "Building VSCode Extension..."
	@cd vscode-extension && npm install && npm run build:webview && npm run compile
	@echo "Extension built"

ext-package: ext-build
	@echo "Packaging VSCode Extension..."
	@cd vscode-extension && npm run package
	@echo "Created vscode-extension/claritask-*.vsix"

ext-install: ext-package
	@echo "Installing VSCode Extension..."
	@code --install-extension vscode-extension/claritask-*.vsix
	@echo "Extension installed"
