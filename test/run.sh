#!/bin/bash

make || exit 1
./bot -balance PLN/BTC -order 50.01PLN -gain 100.5% \
	"$@" ./test/market.py ./test/data.txt
