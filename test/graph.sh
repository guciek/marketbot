#!/bin/bash
# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

O=$(cat grid.txt | grep '^[0-9 ]*$')

GRAPH_COLS=$(($(tput cols) - 20)) || exit 1
TERM_ROWS=$(tput lines) || exit 1
if ! [ $TERM_ROWS -ge 5 -a $GRAPH_COLS -ge 5 ]; then
	echo "Error: invalid terminal size" 1>&2
	exit 1
fi
GRAPH_ROWS=5
while [ $(($GRAPH_ROWS*5)) -lt $(($TERM_ROWS-4)) ]; do
	let GRAPH_ROWS=$GRAPH_ROWS*5
done
while [ $(($GRAPH_ROWS*2)) -lt $(($TERM_ROWS-4)) ]; do
	let GRAPH_ROWS=$GRAPH_ROWS*2
done
let GRAPH_ROWS=$GRAPH_ROWS+1

MAXP=0
while read buy_p sell_p buy_am sell_am; do
	let buy_p=$buy_p*100000
	if [ $MAXP = 0 -o $buy_p -gt $MAXP ]; then
		MAXP=$buy_p
	fi
done < <(echo "$O")

let MAXP=$(($MAXP/50000000+1))*50000000

K=0; while [ $K -lt $GRAPH_ROWS ]; do
	VS[$K]=0
	VB[$K]=0
	let K=$K+1
done

let RSCALE=$MAXP/$(($GRAPH_ROWS-1))

while read buy_p sell_p buy_am sell_am; do
	let buy_p=$buy_p*100000
	let sell_p=$sell_p*100000
	let I=$buy_p/$RSCALE
	let PART=$buy_p%$RSCALE
	let ADDNEXT=$buy_am*$PART/$RSCALE
	VS[$I]=$((${VS[$I]} + $buy_am - $ADDNEXT))
	VS[$(($I+1))]=$((${VS[$(($I+1))]} + $ADDNEXT))
	let X=$sell_am/100*$sell_p/1000000
	let ADDNEXT=$X*$PART/$RSCALE
	VB[$I]=$((${VB[$I]} + $X - $ADDNEXT))
	VB[$(($I+1))]=$((${VB[$(($I+1))]} + $ADDNEXT))
done < <(echo "$O")

MAXV=0
K=0; while [ $K -lt $GRAPH_ROWS ]; do
	V=$((${VS[$K]}+${VB[$K]}))
	if [ $V -gt $MAXV ]; then
		MAXV=$V
	fi
	let K=$K+1
done

let SCALE=$MAXV/$GRAPH_COLS

K=0; while [ $K -lt $GRAPH_ROWS ]; do
	if [ $(($K % 5)) == 0 ]; then
		let P=$K*$RSCALE
		let P=$P/100000
		printf "% 5d.%03d -|" $(($P/1000)) \
			$(($P-$(($(($P/1000))*1000))))
	else
		printf "           |"
	fi
	N=$(($((${VS[$K]}+$(($SCALE/2)))) / $SCALE))
	I=0; while [ $I -lt $N ]; do echo -n "$"; let I=$I+1; done
	N=$(($((${VB[$K]}+$(($SCALE/2)))) / $SCALE))
	I=0; while [ $I -lt $N ]; do echo -n "*"; let I=$I+1; done
	echo
	let K=$K+1
done
