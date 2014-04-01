#!/usr/bin/python3

import time
import sys
import base64
import json
import hmac
import hashlib
from urllib import parse, request

class Cryptsy:
	def __init__(self, api_key, api_secret):
		self.__api_key = str(api_key).strip()
		self.__api_secret = bytes(str(api_secret).strip(), "ascii")
		self.__api_url = "https://www.cryptsy.com/api"
		self.__tradingpairs = None

	def __q(self, fun, args = {}):
		try:
			args['nonce'] = int(time.time())
			args['method'] = fun
			post_data = bytes(parse.urlencode(args), "utf-8")
			sign = hmac.new(self.__api_secret, post_data, hashlib.sha512).hexdigest()
			headers = {
				'Key' : self.__api_key,
				'Sign': sign,
				'User-Agent': "Mozilla/4.0 (compatible; trading bot)"
			}
			url = self.__api_url
			req = request.Request(url, post_data, headers)
			response = request.urlopen(req, None, 5)
			d = str(response.read(), "utf-8")
			try:
				ret = dict(json.loads(d))
			except Exception as e:
				raise Exception("could not parse JSON")
			if "success" not in ret:
				raise Exception("wrong response format")
			if int(ret["success"]) != 1:
				raise Exception(str(ret["error"]))
			if "return" in ret:
				return ret["return"]
			return dict()
		except Exception as e:
			raise Exception("API request failed: "+str(e))

	def getTradingPairs(self):
		if self.__tradingpairs == None:
			r = self.__q("getmarkets")
			pairs = dict()
			self.__pairbyid = dict()
			for m in r:
				pairname = (str(m["primary_currency_code"])+"_"+
					str(m["secondary_currency_code"])).upper()
				mid = str(m["marketid"])
				pairs[pairname] = mid
				self.__pairbyid[mid] = pairname
			self.__tradingpairs_set = set(pairs.keys())
			self.__tradingpairs = pairs
		return self.__tradingpairs_set

	def getFunds(self):
		return self.__q("getinfo")

	def getOrders(self):
		self.getTradingPairs()
		r = self.__q("allmyorders")
		for o in r:
			o["pair"] = self.__pairbyid[str(o["marketid"])]
		return r

	def placeOrder(self, pair, tpe, amount, price):
		self.getTradingPairs()
		marketid = self.__tradingpairs[pair]
		args = {"marketid": str(marketid), "ordertype": str(tpe),
			"quantity": str(amount), "price": ("%.15f"%float(price))}
		return self.__q("createorder", args)

	def cancelOrder(self, oid):
		return self.__q("cancelorder", {"orderid": str(oid)})

def cmdLine(market, line):
	if line == "": return False
	line = line.split(" ")

	if line[0] == "time":
		print("time "+str(int(time.time())))
		return True

	if line[0] == "echo":
		print("echo "+line[1])
		return True

	if line[0] == "wait":
		try:
			time.sleep(10)
		except:
			pass
		print("ok wait")
		return True

	if (line[0] == "buy") and (line[3] == "for"):
		am1 = line[1]
		cur1 = line[2].upper()
		am2 = line[4]
		cur2 = line[5].upper()
		trading_pairs = market.getTradingPairs()
		tpe = "Buy"
		if cur2+"_"+cur1 in trading_pairs:
			am1, am2 = am2, am1
			cur1, cur2 = cur2, cur1
			tpe = "Sell"
		if cur1+"_"+cur2 in trading_pairs:
			pr = float(am2)/float(am1)
			market.placeOrder(cur1+"_"+cur2, tpe, am1, pr)
			print("ok buy")
		else:
			print(trading_pairs)
			print("error unsupported trading pair '%s'"%(cur1+"_"+cur2))
		return True

	if line[0] == "cancel":
		r = market.cancelOrder(str(line[1]))
		if "error" in r:
			print("error", r["error"])
		else:
			print("ok cancel")
		return True

	if line[0] == "orders":
		r = market.getOrders()
		print("orders:")
		for o in r:
			oid = str(o["orderid"])
			currencies = str(o["pair"]).split("_")
			amounts = [str(o["quantity"]), str(o["total"])]
			if str(o["ordertype"]) == "Buy":
				print("%s buy %s %s for %s %s" % (
					oid,
					amounts[0], currencies[0],
					amounts[1], currencies[1]
				))
			elif str(o["ordertype"]) == "Sell":
				print("%s buy %s %s for %s %s" % (
					oid,
					amounts[1], currencies[1],
					amounts[0], currencies[0]
				))
			else:
				print("unknown order type ", o["ordertype"])
		print(".")
		return True

	elif line[0] == "totalbalance":
		r = market.getFunds()
		print("totalbalance:")
		if "balances_available" in r:
			for k in r["balances_available"]:
				val = str(r["balances_available"][k])
				cur = str(k)
				if cur.isalpha() and (len(cur) <= 10):
					print(val+" "+cur)
		if "balances_hold" in r:
			for k in r["balances_hold"]:
				val = str(r["balances_hold"][k])
				cur = str(k)
				if cur.isalpha() and (len(cur) <= 10):
					print(val+" "+cur)
		print(".")
		return True

	if line[0] == "exit":
		print("exit")
		return False

	print("error Unknown command")
	return True

def run():
	b = Cryptsy(sys.argv[1], sys.argv[2])
	for line in sys.stdin:
		try:
			if not cmdLine(b, line.strip()):
				break
			sys.stdout.flush()
			try:
				time.sleep(1)
			except:
				pass
		except IOError:
			raise
		except KeyboardInterrupt:
			raise
		except Exception as e:
			print("error", e)
			sys.stdout.flush()

try:
	run()
except IOError:
	pass
except KeyboardInterrupt:
	pass
except Exception as e:
	sys.stderr.write("Error: "+str(e)+"\n")
	sys.stderr.flush()
