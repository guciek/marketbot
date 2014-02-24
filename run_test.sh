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

echo "buy_p sell_p buy_am sell_am" > grid.txt
P=2000; while [ $P -lt 3000 ]; do
	let SELLP=$P+$(($P*3/100))
	if [ $P -ge 2450 ]; then
		BUY=0
		SELL=$((10000000000/$SELLP))
	else
		BUY=10000000
		SELL=0
	fi
	echo $P $SELLP $BUY $SELL
	let P=$P+$(($P/200))
done >> grid.txt

./bot ./test-market.py test-data.log || exit 1

exit 0
