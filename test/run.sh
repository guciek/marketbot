#!/bin/bash

make || exit 1
./bot -balance PLN/BTC -order 50.01PLN -fee 0.4% \
	"$@" ./test/market.py ./test/data.txt
