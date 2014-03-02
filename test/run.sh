#!/bin/bash

make || exit 1
./bot -natural 2500PLN/BTC -order 50.01PLN -gain 101.6% \
	"$@" ./test/market.py ./test/data.txt
