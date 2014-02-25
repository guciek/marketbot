#!/bin/bash

make || exit 1
./bot -natural '$51' -target 2500 -fee 996 "$@" ./test/market.py ./test/data.txt
