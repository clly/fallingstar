
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
docker/build:
	@docker build -t $(APPNAME) .
