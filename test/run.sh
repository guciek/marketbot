#!/bin/bash

make || exit 1
./bot \
	-fee 0.4% \
	-natural 2500pln/BTC,160pln \
	"$@" ./test/market.py ./test/data.txt
