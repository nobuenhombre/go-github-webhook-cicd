include .env
export $(shell sed 's/=.*//' .env)
#======================================================
# SERVER_ROLE read from env
# values
#   - production
#   - develop
#   - develop.local
#   - etc...
# see also /configs/${SERVER_ROLE} dirs with configs
#======================================================

#======================================================
# HTTP Server
#======================================================
PROJECT=go-github-webhook-cicd
APP_VERSION=v0.0.1
APP_NAME=go-github-webhook-cicd
APP_LOG_NAME=go-github-webhook-cicd

PROJECT_ROOT_PATH=../../../..
INSTALL_PATH=/usr/local/bin

BUILD_PLATFORM=linux
BIN_PATH=bin/$(APP_NAME)/linux
APP_BINARY=$(BIN_PATH)/$(APP_NAME)

SERVICE_NAME="api_go_github_webhook_cicd"
SERVICE_PATH=/etc/systemd/system

#=========================================================================
#
#=========================================================================
.PHONY: help build

help: Makefile
	@echo "Выберите опцию сборки "$(BINARY_NAME)":"
	@sed -n 's/^##//p' $< | column -s ':' |  sed -e 's/^/ /'

## all: Compile the application and install
all: service-stop uninstall-service uninstall-app build-app install-app install-service service-start

## build-app: Compile the application
build-app:
	cd $(PROJECT_ROOT_PATH)/ && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=$(BUILD_PLATFORM) GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_BINARY) -v src/cmd/$(APP_NAME)/main.go && \
	chmod +x $(APP_BINARY);

## install-app: Install the app
install-app:
	sudo mkdir -p /var/log/$(APP_LOG_NAME) && \
    sudo chmod 777 /var/log/$(APP_LOG_NAME) && \
    sudo ln -sf $(shell pwd)/$(PROJECT_ROOT_PATH)/$(APP_BINARY) $(INSTALL_PATH)/$(APP_NAME);

## uninstall-app: Delete the application
uninstall-app:
	sudo rm -f $(INSTALL_PATH)/$(APP_NAME);

## install-service: Install the HTTP Service
install-service:
	sudo systemctl enable $(shell pwd)/$(PROJECT_ROOT_PATH)/configs/$(SERVER_ROLE)/$(SERVICE_NAME).service && \
	sudo systemctl daemon-reload;

## uninstall-service: Delete the HTTP Service
uninstall-service:
	sudo systemctl disable $(shell pwd)/$(PROJECT_ROOT_PATH)/configs/$(SERVER_ROLE)/$(SERVICE_NAME).service && \
	sudo systemctl daemon-reload;

## service-stop: Stop the HTTP Service
service-stop:
	sudo systemctl stop $(SERVICE_NAME);

## service-start: Start the HTTP Service
service-start:
	sudo systemctl start $(SERVICE_NAME);

## service-status: Get Status the HTTP Service
service-status:
	sudo systemctl status $(SERVICE_NAME);

## service-restart: Restart the HTTP Service
service-restart:
	sudo systemctl restart $(SERVICE_NAME);

## upgrade: Upgrade from sources
upgrade: service-stop build-app service-start service-status
