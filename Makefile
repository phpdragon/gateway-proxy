include .env

##############################################
## Go项目编译脚本
## 参考：https://studygolang.com/articles/14919?fr=sidebar
##############################################

PROJECT_NAME=$(shell basename "$(PWD)")

# Go related variables.
CURRENT_DIR=$(shell pwd)

#https://www.cnblogs.com/blue-sea-sky/p/5689181.html
GO_HOME_DIR=$(shell go env | grep 'GOROOT' | awk -F '=' '{print $$2}' | sed 's/"//g')
GO_CMD=$(GO_HOME_DIR)/bin/go

APP_BIN_DIR=$(CURRENT_DIR)/bin
GO_FILES=$(wildcard src/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## deps: Install missing dependencies.
deps:
	@echo "Checking if there is any missing dependencies..."
	$(GO_CMD) get github.com/astaxie/beego/orm
	$(GO_CMD) get github.com/go-redis/redis
	$(GO_CMD) get github.com/go-sql-driver/mysql
	$(GO_CMD) get github.com/go-yaml/yaml
	$(GO_CMD) get github.com/natefinch/lumberjack
	$(GO_CMD) get github.com/streadway/amqp
	$(GO_CMD) get go.uber.org/zap
	$(GO_CMD) get go.uber.org/zap/zapcore
	@echo "Checking any missing dependencies is over!"
	@echo

## build: Compile the binary.
build: clean
	@echo "Building binary..."
	#go build 参数说明 https://www.cnblogs.com/davygeek/p/6386035.html
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(APP_BIN_DIR)/$(PROJECT_NAME) -a -x -v $(GO_FILES)
	@-chmod a+x -R $(APP_BIN_DIR)/*
	@echo "Compile the complete!"
	@echo

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo "Cleaning build cache"
	@$(GO_CMD) clean
	@-rm -f $(APP_BIN_DIR)/$(PROJECT_NAME)
	@echo "Clean build cache over!"
	@echo

## package: Package the app
package: build
	@echo "Taring project package..."
	@-tar -czvf $(CURRENT_DIR)/$(PROJECT_NAME).tar.gz -C $(CURRENT_DIR) bin etc log favicon.ico LICENSE README*.md
	@echo "Taring project package over!"
	@echo

## deploy: Deploy package to server site
deploy: package
	scp $(CURRENT_DIR)/$(PROJECT_NAME).tar.gz root@192.168.1.2:/data/server/

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
