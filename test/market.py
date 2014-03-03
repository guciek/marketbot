#!/usr/bin/python3
# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

import time
import sys
import json
from urllib import request
from random import randint

class FakeMarket:
	def __init__(self):
		self.__cash = 10000.0
		self.__btc = 0.0
		self.__cashout_price = -1.0
		self.__fee = 0.996
		self.__orders = []
		self.__nextoid = 1000
		self.__lastprice = -1.0

	def __runTransaction(self):
		for i in range(0, len(self.__orders)):
			o = self.__orders[i]
			if o["type"] == "buy":
				if o["price"] >= self.__lastprice:
					self.__btc = self.__btc + o["amount"]*self.__fee
					del self.__orders[i]
					return True
			if o["type"] == "sell":
				if o["price"] <= self.__lastprice:
					self.__cash = self.__cash + o["amount"]*o["price"]*self.__fee
					del self.__orders[i]
					return True
		return False

	def onPriceChange(self, pr):
		if self.__cashout_price < 0.0:
			self.__cash = self.__cash/2
			self.__cashout_price = pr
			self.__btc = self.__cash/(pr*self.__fee)
		self.__lastprice = pr
		while self.__runTransaction():
			pass

	def info(self):
		f = self.getFunds()["funds"]
		return (
			("%0.5f BTC + %0.2f PLN, "+
				"cash out %0.2f PLN at %0.2f PLN/BTC") %
			(
				f["BTC"], f["PLN"],
				f["PLN"] + (f["BTC"] * self.__cashout_price * self.__fee),
				self.__cashout_price
			)
		)

	def getOrders(self):
		return {"orders": self.__orders}

	def getFunds(self):
		sum_cash = self.__cash
		sum_btc = self.__btc
		for o in self.__orders:
			if o["type"] == "sell":
				sum_btc = sum_btc + o["amount"]
			if o["type"] == "buy":
				sum_cash = sum_cash + o["amount"]*o["price"]
		return {"funds": {"PLN": sum_cash, "BTC": sum_btc}}

	def cancelOrder(self, oid):
		for i in range(0, len(self.__orders)):
			o = self.__orders[i]
			if str(o["id"]) == str(oid):
				if o["type"] == "sell":
					self.__btc = self.__btc + o["amount"]
				if o["type"] == "buy":
					self.__cash = self.__cash + o["amount"]*o["price"]
				del self.__orders[i]
				return {"result": "ok"}
		return {"error": "order id not found"}

	def buySell(self, op, amount_btc, price):
		amount_btc = float(amount_btc)*(1.0+randint(-10,10)*0.000001)
		price = float(price)*(1.0+randint(-10,10)*0.000001)
		if price*amount_btc < 50.0:
			return {"error": "order too small"}
		if op == "buy":
			if self.__cash < amount_btc*price:
				return {"error": "not enough funds"}
			self.__cash = self.__cash - amount_btc*price
		elif op == "sell":
			if self.__btc < amount_btc:
				return {"error": "not enough funds"}
			self.__btc = self.__btc - amount_btc
		else:
			return {"error": "invalid order type"}
		self.__orders.append({"type": op, "amount": amount_btc,
			"price": price, "id": "#"+str(self.__nextoid)})
		self.__nextoid = self.__nextoid+1
		while self.__runTransaction():
			pass
		return {"result": "ok"}

def fakeTime(fn, market, store = dict()):
	store["f"] = open(fn, "rb")
	l = str(store["f"].readline(), "ascii").split(" ")
	store["t"] = int(l[0])-1
	store["next_ts"] = int(l[0])
	store["next_price"] = int(l[1])
	def passTime(step):
		store["t"] = store["t"]+step
		while store["t"] >= store["next_ts"]:
			market.onPriceChange(store["next_price"]*0.00001)
			l = str(store["f"].readline(), "ascii")
			if l == "":
				store["f"].close()
				return False
			l = l.split(" ")
			store["next_ts"] = int(l[0])
			store["next_price"] = int(l[1])
		return True
	def getTime():
		return store["t"]
	passTime(1)
	return passTime, getTime

def cmdLine(market, line, passTime, getTime):
	if line == "": return False
	line = line.split(" ")

	if line[0] == "time":
		print("time "+str(int(getTime())))
		return True

	if line[0] == "echo":
		print("echo "+line[1])
		return True

	if line[0] == "wait":
		passTime(60)
		print("ok wait")
		return True

	if randint(1, 100) <= 3:
		print("error")
		return True

	if (line[0] == "buy") and (line[3] == "for"):
		am1 = line[1]
		am2 = line[4]
		if (line[2] == "PLN") and (line[5] == "BTC"):
			r = market.buySell("sell", am2, float(am1)/float(am2))
			if "error" in r:
				print("error", r["error"])
			else:
				print("ok buy")
		elif (line[2] == "BTC") and (line[5] == "PLN"):
			r = market.buySell("buy", am1, float(am2)/float(am1))
			if "error" in r:
				print("error", r["error"])
			else:
				print("ok buy")
		else:
			print("error unsupported trading pair")
		return True

	if line[0] == "cancel":
		r = market.cancelOrder(line[1])
		if "error" in r:
			print("error", r["error"])
		else:
			print("ok cancel")
		return True

	if line[0] == "orders":
		r = market.getOrders()
		if "orders" not in r:
			print("error")
			return True
		print("orders:")
		for o in r["orders"]:
			t = str(o["type"])
			pr = float(o["price"])
			am_btc = float(o["amount"])
			am_pln = am_btc*pr
			oid = str(o["id"])
			if t == "sell":
				print("%s buy %.5f PLN for %.8f BTC" % (oid, am_pln, am_btc))
			elif t == "buy":
				print("%s buy %.8f BTC for %.5f PLN" % (oid, am_btc, am_pln))
			else:
				print("unknown order type")
		print(".")
		return True

	elif line[0] == "totalbalance":
		r = market.getFunds()
		if "funds" not in r:
			print("error")
			return True
		print("totalbalance:")
		for k in r["funds"]:
			print("%s %s" % (str(r["funds"][k]), str(k).upper()))
		print(".")
		return True

	if line[0] == "exit":
		print("exit")
		return False

	print("error Unknown command")
	return True

def run():
	b = FakeMarket()
	passTime, getTime = fakeTime(sys.argv[1], b)
	try:
		for line in sys.stdin:
			if not cmdLine(b, line.strip(), passTime, getTime):
				break
			sys.stdout.flush()
			passTime(1)
	finally:
		sys.stderr.write("[Market] "+b.info()+"\n")

try:
	run()
except IOError:
	pass
except KeyboardInterrupt:
	pass
except Exception as e:
	sys.stderr.write("[Market] Error: "+str(e)+"\n")
