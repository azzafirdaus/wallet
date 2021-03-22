#!/bin/bash

build:
	@go build -v

run:
	@echo "CONFIGURING YOUR MACHINE FOR DEVELOPMENT ⚙️ ⚙️ ⚙️ "
	@go mod vendor -v
	@go build -v && ./wallet