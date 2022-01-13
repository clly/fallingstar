MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := all
APPNAME := "goregex"

BIN_DIR ?= $(shell go env GOPATH)/bin
export PATH := $(PATH):$(BIN_DIR)

ifneq ("$(wildcard .makefiles/*.mk)","")
	include .makefiles/*.mk
else
    $(info "no makefiles to load")
endif

.PHONY: all
all: test build

