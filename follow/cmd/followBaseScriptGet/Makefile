SHELL      = /bin/bash

Target     = $(shell basename $(abspath $(dir $$PWD)))
Prefix     = ../../../../release

all:
	source /opt/toolchain-sunxi/environment-setup-arm-openwrt-linux && \
	GO111MODULE=off GOARM=7 GOARCH=arm GOOS=linux CGO_ENABLED=1 \
	go build -v -o $(Prefix)/$(Target).test
