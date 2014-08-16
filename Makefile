# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

export GOPATH := $(shell pwd)
SRCFILES := $(shell find src -type f)
MODULES := $(shell ls src)

bot: Makefile $(SRCFILES)
	@go test $(MODULES)
	@go build $@
	@strip $@
