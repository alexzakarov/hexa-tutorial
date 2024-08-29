#!make
SERVICE_NAME = mono
DOCKER_NAME = core-api

include $(CURDIR)/config/dev.env
include $(CURDIR)/init/service.mak
include $(CURDIR)/init/nginx.mak

PATH_SERVICE_FILE := /etc/systemd/system/${APP_SERVICE_NAME}.service
PATH_NGINX_MAIN_CONF_FILE := /etc/nginx/conf.d/apis.conf
PATH_NGINX_MAIN_CONF_CONTENT := "include /etc/nginx/api_conf.d/*.conf;"
PATH_NGINX_DIR := /etc/nginx/api_conf.d
PATH_NGINX_FILE := ${PATH_NGINX_DIR}/${APP_ROUTE}.${APP_DOMAIN_MAIN}.conf
SERVICE_STATUS := $(shell systemctl is-enabled ${APP_SERVICE_NAME}.service)

start:
	make swag
	[ ${SERVICE_STATUS} = "enabled" ] &&  make rebuild || make install

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/$(SERVICE_NAME)/main.go

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	go build -ldflags "-s -w" -o /${APP_DIR}/${APP_SERVICE_NAME}/${APP_EXE_NAME} ./cmd/${APP_SERVICE_PATH}/main.go
	mkdir -p /${APP_DIR}/${APP_SERVICE_NAME}/config
	mkdir -p /${APP_DIR}/${APP_SERVICE_NAME}/doc
	mkdir -p /${APP_DIR}/${APP_SERVICE_NAME}/i18n
	cp -r config/* /${APP_DIR}/${APP_SERVICE_NAME}/config
	cp -r docs/* /${APP_DIR}/${APP_SERVICE_NAME}/doc
	cp -r i18n/* /${APP_DIR}/${APP_SERVICE_NAME}/i18n
	[ ! -f ${PATH_NGINX_MAIN_CONF_FILE} ] && echo ${PATH_NGINX_MAIN_CONF_CONTENT} > ${PATH_NGINX_MAIN_CONF_FILE} || echo "${PATH_NGINX_MAIN_CONF_FILE} is exist"
	[ ! -d ${PATH_NGINX_DIR} ] && mkdir ${PATH_NGINX_DIR}; echo "$$nginx" > ${PATH_NGINX_FILE} || echo "$$nginx" > ${PATH_NGINX_FILE}
	nginx -s reload
	[ ! -f ${PATH_SERVICE_FILE} ] && echo "$$service" > ${PATH_SERVICE_FILE} || rm -rf ${PATH_SERVICE_FILE}; echo "$$service" > ${PATH_SERVICE_FILE}
	systemctl daemon-reload
	systemctl enable ${APP_SERVICE_NAME}.service
	systemctl restart ${APP_SERVICE_NAME}.service

uninstall:
	rm -rf ${PATH_NGINX_FILE}
	nginx -s reload
	systemctl disable ${APP_SERVICE_NAME}.service
	systemctl stop ${APP_SERVICE_NAME}.service
	rm -rf /${APP_DIR}/${APP_SERVICE_NAME}
	rm -rf ${PATH_SERVICE_FILE}
	systemctl daemon-reload

rebuild:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	go build -ldflags "-s -w" -o /${APP_DIR}/${APP_SERVICE_NAME}/${APP_EXE_NAME} ./cmd/${APP_SERVICE_PATH}/main.go
	mkdir -p /${APP_DIR}/${APP_SERVICE_NAME}/doc
	cp -r docs/* /${APP_DIR}/${APP_SERVICE_NAME}/doc
	cp -r config/* /${APP_DIR}/${APP_SERVICE_NAME}/config
	cp -r i18n/* /${APP_DIR}/${APP_SERVICE_NAME}/i18n
	systemctl restart ${APP_SERVICE_NAME}.service

ports:
	sudo netstat -tulpn | grep LISTEN