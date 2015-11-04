#!/usr/bin/python3

import time
import sys
import base64
import json
import hmac
import hashlib
import urllib
from urllib import parse, request

class Bitfinex:
    def __init__(self, api_key, api_secret):
        self.__api_key = str(api_key).strip()
        self.__api_secret = bytes(str(api_secret).strip(), "ascii")
        self.__tradingpairs = None

    def __q(self, fun, args = {}):
        try:
            args['request'] = "/v1/" + fun
            args['nonce'] = str(int(time.time() * 10000))
            post_data = base64.standard_b64encode(bytes(json.dumps(args), "ascii"))
            sign = hmac.new(self.__api_secret, post_data, hashlib.sha384).hexdigest()
            headers = {
                'Content-Type': "application/json",
                'Accept': "application/json",
                'X-BFX-APIKEY': self.__api_key,
                'X-BFX-PAYLOAD': str(post_data, "ascii"),
                'X-BFX-SIGNATURE': sign
            }
            url = "https://api.bitfinex.com" + args['request']
            req = request.Request(url, bytes(), headers)
            response = request.urlopen(req, None, 5)
            d = str(response.read(), "ascii")
            ret = {}
            try:
                ret = json.loads(d)
            except Exception as e:
                raise Exception("could not parse JSON")
            return ret
        except urllib.error.HTTPError as e:
            msg = e
            try:
                c = json.loads(str(e.read(), "ascii"))["message"]
                msg = c
            except Exception as ee:
                pass
            raise Exception("API request failed: "+str(msg))
        except Exception as e:
            raise Exception("API request failed: "+str(e))

    def getTradingPairs(self):
        return ["BTCUSD"]

    def getFunds(self):
        return self.__q("balances")

    def getOrders(self):
        r = self.__q("orders")
        return r

    def placeOrder(self, pair, tpe, amount, price):
        args = {"symbol": str(pair).lower(), "amount": str(amount),
            "price": ("%.4f"%float(price)), "exchange" : "bitfinex",
            "side": str(tpe), "type": "exchange limit"
            }
        self.__q("order/new", args)

    def cancelOrder(self, oid):
        return self.__q("order/cancel", {"order_id": int(oid)})

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
        if cur2+cur1 in trading_pairs:
            am1, am2 = am2, am1
            cur1, cur2 = cur2, cur1
            tpe = "sell"
        if cur1+cur2 in trading_pairs:
            pr = float(am2)/float(am1)
            market.placeOrder(cur1+cur2, tpe, am1, pr)
            print("ok buy")
        else:
            print("error unsupported trading pair '%s'"%(cur1+cur2))
        return True

    if line[0] == "cancel":
        market.cancelOrder(str(line[1]))
        print("ok cancel")
        return True

    if line[0] == "orders":
        r = market.getOrders()
        print("orders:")
        for o in r:
            if str(o["is_cancelled"]).lower() == "false":
                if str(o["type"]).startswith("exchange "):
                    oid = str(o["id"])
                    currencies = str(o["symbol"])
                    currencies = [currencies[0:3], currencies[3:6]]
                    amounts = [str(o["remaining_amount"]), str(float(o["remaining_amount"]) * float(o["price"]))]
                    if str(o["side"]) == "buy":
                        print("%s buy %s %s for %s %s" % (
                            oid,
                            amounts[0], currencies[0],
                            amounts[1], currencies[1]
                        ))
                    elif str(o["side"]) == "sell":
                        print("%s buy %s %s for %s %s" % (
                            oid,
                            amounts[1], currencies[1],
                            amounts[0], currencies[0]
                        ))
                    else:
                        print("unknown order type ", o["side"])
            elif str(o["is_cancelled"]).lower() == "true":
                pass
            else:
                print("unknown order status ", o["is_cancelled"])
        print(".")
        return True

    elif line[0] == "totalbalance":
        r = market.getFunds()
        print("totalbalance:")
        for k in r:
            if k["type"] == "exchange":
                print(k["amount"] + " " + k["currency"])
        print(".")
        return True

    if line[0] == "exit":
        print("exit")
        return False

    print("error Unknown command")
    return True

def run():
    b = Bitfinex(sys.argv[1], sys.argv[2])
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
