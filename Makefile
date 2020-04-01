include .env

PROJECT_NAME=$(shell basename "$(PWD)")

# Go related variables.
CURRENT_DIR=$(shell pwd)

#https://www.cnblogs.com/blue-sea-sky/p/5689181.html
GO_CMD=$(shell whereis go | awk '{print $$2}')

APP_BIN_DIR=$(CURRENT_DIR)/bin
GO_FILES=$(wildcard src/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## deps: Install missing dependencies.
deps:
	@echo "  >  Checking if there is any missing dependencies..."
#	$(GO_CMD) get github.com/markbates/goth
#	$(GO_CMD) get github.com/markbates/pop

## build: Compile the binary.
build: clean deps
	@echo "  >  Building binary..."
	#go build 参数说明 https://www.cnblogs.com/davygeek/p/6386035.html
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(APP_BIN_DIR)/$(PROJECT_NAME) -a -x -v $(GO_FILES)
	@echo "  >  Compile the complete!"

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo "  >  Cleaning build cache"
	@$(GO_CMD) clean
	@-rm -f $(APP_BIN_DIR)/$(PROJECT_NAME)

## package: Package the app
package: build
	@-tar -czvf $(CURRENT_DIR)/$(PROJECT_NAME).tar.gz -C $(CURRENT_DIR) bin etc log favicon.ico LICENSE README*.md

## deply: Deply package to server site
deply:
	#scp db_svr im@192.168.251.53:/home/im/bin/db/

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
