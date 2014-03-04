#!/bin/bash

make || exit 1
./bot -natural 2500PLN/BTC -order 50.01PLN -fee 0.4% -gain 0.5% -pricegain 10PLN/BTC \
	"$@" ./test/market.py ./test/data.txt
