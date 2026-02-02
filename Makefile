.PHONY: install uninstall

BINARY_NAME := clari
INSTALL_PATH := $(if $(GOBIN),$(GOBIN),$(if $(GOPATH),$(GOPATH)/bin,$(HOME)/go/bin))

install:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/claritask
	@mkdir -p $(INSTALL_PATH)
	@mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_PATH)/$(BINARY_NAME)"

uninstall:
	@rm -f $(BINARY_NAME)
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"
