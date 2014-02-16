#!/bin/bash
# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

make || exit 1

TMPDIR=$(mktemp -d) || exit 1
trap "rm -rf $TMPDIR" EXIT

cp bot test/* "$TMPDIR" || exit 1
cd "$TMPDIR" || exit 1
./bot ./test-market.py test-data.log || exit 1

exit 0
