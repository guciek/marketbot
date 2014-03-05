#!/bin/bash

make || exit 1
./bot \
	-fee 0.4% \
	-natural 2500pln/BTC,50.01pln \
	-mingain 1%,20pln/btc \
	"$@" ./test/market.py ./test/data.txt
