#!/bin/bash

make || exit 1
./bot -natural '$50.01' -target 2.5 -fee 0.4% -spread 2.5% "$@" ./test/market.py ./test/data.txt
