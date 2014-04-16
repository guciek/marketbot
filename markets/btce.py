#!/usr/bin/python3

import time
import sys
import base64
import json
import hmac
import hashlib
from urllib import parse, request

class Btce:
	def __init__(self, api_key, api_secret):
		self.__api_key = str(api_key).strip()
		self.__api_secret = bytes(str(api_secret).strip(), "ascii")
		self.__api_url = "https://btc-e.com/tapi"

	def __q(self, fun, args = {}):
		try:
			try:
				time.sleep(1)
			except:
				pass
			args['nonce'] = int(time.time())
			args['method'] = fun
			post_data = bytes(parse.urlencode(args), "utf-8")
			sign = hmac.new(self.__api_secret, post_data, hashlib.sha512).hexdigest()
			headers = {
				'Content-type': "application/x-www-form-urlencoded",
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
				if (fun == "ActiveOrders") and (str(ret["error"]) == "no orders"):
					return {}
				raise Exception(str(ret["error"]))
			if "return" in ret:
				return ret["return"]
			return dict()
		except Exception as e:
			raise Exception("API request failed: "+str(e))

	def getTradingPairs(self):
		return set(["BTC_USD"])

	def getInfo(self):
		return self.__q("getInfo")

	def getOrders(self):
		r = self.__q("ActiveOrders")
		for o in r:
			r[o]["pair"] = r[o]["pair"].upper()
		return r

	def placeOrder(self, pair, tpe, amount, price):
		args = {"pair": str(pair).lower(), "type": str(tpe),
			"amount": ("%.8f"%float(amount)), "rate": ("%.2f"%float(price))}
		return self.__q("Trade", args)

	def cancelOrder(self, oid):
		return self.__q("CancelOrder", {"order_id": str(oid)})

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
		tpe = "buy"
		if cur2+"_"+cur1 in trading_pairs:
			am1, am2 = am2, am1
			cur1, cur2 = cur2, cur1
			tpe = "sell"
		if cur1+"_"+cur2 in trading_pairs:
			pr = float(am2)/float(am1)
			market.placeOrder(cur1+"_"+cur2, tpe, am1, pr)
			print("ok buy")
		else:
			print("error unsupported trading pair '%s'"%(cur1+"_"+cur2))
		return True

	if line[0] == "cancel":
		market.cancelOrder(str(line[1]))
		print("ok cancel")
		return True

	if line[0] == "orders":
		r = market.getOrders()
		print("orders:")
		for oid in r:
			o = r[oid]
			currencies = str(o["pair"]).split("_")
			amounts = [str(o["amount"]),
				str(float(o["amount"])*float(o["rate"]))]
			if str(o["type"]) == "buy":
				print("%s buy %s %s for %s %s" % (
					oid,
					amounts[0], currencies[0],
					amounts[1], currencies[1]
				))
			elif str(o["type"]) == "sell":
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
		r1 = market.getInfo()
		r2 = market.getOrders()
		print("totalbalance:")
		for k in r1["funds"]:
			val = str(r1["funds"][k])
			cur = str(k)
			if cur.isalpha() and (len(cur) <= 10) and (val != "0"):
				print(val+" "+cur.upper())
		for oid in r2:
			o = r2[oid]
			currencies = str(o["pair"]).split("_")
			amounts = [str(o["amount"]),
				str(float(o["amount"])*float(o["rate"]))]
			if str(o["type"]) == "buy":
				print(amounts[1]+" "+currencies[1])
			elif str(o["type"]) == "sell":
				print(amounts[0]+" "+currencies[0])
			else:
				print("unknown order type ", o["ordertype"])
		print(".")
		return True

	if line[0] == "exit":
		print("exit")
		return False

	print("error Unknown command")
	return True

def run():
	b = Btce(sys.argv[1], sys.argv[2])
	for line in sys.stdin:
		try:
			if not cmdLine(b, line.strip()):
				break
			sys.stdout.flush()
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
