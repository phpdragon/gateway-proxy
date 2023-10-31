include .env

PROJECT_NAME=$(shell basename "$(PWD)")

# Go related variables.
CURRENT_DIR=$(shell pwd)

GO_HOME_DIR=$(shell go env | grep 'GOROOT' | awk -F '=' '{print $$2}' | sed 's/"//g')
GO_CMD=$(GO_HOME_DIR)/bin/go

APP_BIN_FILE=$(CURRENT_DIR)/$(PROJECT_NAME)
APP_SCRIPTS_DIR=$(CURRENT_DIR)/scripts
APP_GO_FILES=$(wildcard cmd/$(PROJECT_NAME)/*.go)
APP_TAR_FILE=$(CURRENT_DIR)/$(PROJECT_NAME).tar.gz
APP_DEPLOY_DIR=/data/server/$(PROJECT_NAME)/
APP_DEPLOY_ACCOUNT="root@192.168.1.2"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# 去掉 DWARF 调试信息，得到的程序就不能用 gdb 调试了
LDFLAGS="-w"

## deps: Install missing dependencies(安装丢失的依赖项).
deps:
	@echo "Checking if there is any missing dependencies..."
	$(GO_CMD) env -w GO111MODULE=on
	$(GO_CMD) env -w GOPROXY=https://goproxy.cn,direct
	$(GO_CMD) mod download
	@echo "Checking any missing dependencies is over!"
	@echo

## clean: Clean build files. Runs `go clean` internally(清理构建文件).
clean:
	@echo "Cleaning build cache"
	@$(GO_CMD) clean
	@-rm -f $(APP_BIN_FILE)
	@-rm -f $(CURRENT_DIR)/$(PROJECT_NAME).tar.gz
	@echo "Clean build cache over!"
	@echo

## build: Compile the binary(编译二进制文件).
build: clean
	@echo "Building binary..."
	#go build 参数说明 https://www.cnblogs.com/davygeek/p/6386035.html
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(APP_BIN_FILE) -ldflags "$(LDFLAGS)" -a -x -v $(APP_GO_FILES)
	@-chmod a+x -R $(APP_BIN_FILE) $(APP_SCRIPTS_DIR)/*.sh
	@-dos2unix $(APP_SCRIPTS_DIR)/*.sh
	@echo "Compile the complete!"
	@echo

## package: Package the app(打包应用程序).
package: build
	@echo "Taring project package..."
	@-tar -czvf $(APP_TAR_FILE) -C $(CURRENT_DIR) configs logs $(PROJECT_NAME) scripts/server.sh favicon.ico --exclude=*.log
	@echo "Taring project package over!"
	@echo

## deploy: Deploy package to server site(将包部署到服务器站点).
deploy: package
	@-sshpass -p password scp $(APP_TAR_FILE) $(APP_DEPLOY_ACCOUNT):$(APP_DEPLOY_DIR)
	@-sshpass -p password $(APP_DEPLOY_ACCOUNT) "cd $(APP_DEPLOY_DIR) && tar zvf $(PROJECT_NAME).tar.gz && unlink $(PROJECT_NAME).tar.gz && sh scripts/server.sh restart"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " #############################################################"
	@echo " #                       Go项目编译脚本                      #"
	@echo " #  参考：https://studygolang.com/articles/14919?fr=sidebar  #"
	@echo " #############################################################"
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)"(请选择运行的命令):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
