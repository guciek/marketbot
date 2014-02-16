# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

GOFILES := $(wildcard src/*.go)
SRCFILES := $(patsubst %_test.go,,$(GOFILES))

bot: Makefile $(GOFILES)
	@go test $(GOFILES)
	@go build -o $@ $(SRCFILES)
