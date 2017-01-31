package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"gopkg.in/redis.v4"
)

type tradeAPI interface {
	pushOrder() error
}

type bitfinexAPI struct {
	key    string
	secret string
}

// creates a market order using the key and secret you hopefully provided
// only market orders and btcusd for now
func (b *bitfinexAPI) marketOrder(amount, price int) error {

	req, err := http.NewRequest("POST", "https://api.bitfinex.com/v1/order/new", nil)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"request": "/v1/order/new",
		"nonce":   fmt.Sprintf("%v", time.Now().Unix()*10000),
		"symbol":  "btcusd",
		"amount":  amount,
		"price":   price,
		"type":    "market",
	}

	// if only there was a convenient protocol that could store information in a
	// packet of sorts, ill just json encode the information and call
	// it a day
	payload_json, _ := json.Marshal(payload)
	payload_enc := base64.StdEncoding.EncodeToString(payload_json)

	// how about we also make em run through an encryption scheme???
	sig := hmac.New(sha512.New384, []byte(b.secret))
	sig.Write([]byte(payload_enc))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-BFX-APIKEY", b.key)
	req.Header.Add("X-BFX-PAYLOAD", payload_enc)
	req.Header.Add("X-BFX-SIGNATURE", hex.EncodeToString(sig.Sum(nil)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

const bitfinexWS = "wss://api2.bitfinex.com:3000/ws/2"

func connectBitfinex() {

	monitorWebsocket(
		bitfinexWS,
		getPayload("book", "BTCUSD"),
		onBitfinexBookMessage("btcusd"))

	monitorWebsocket(
		bitfinexWS,
		getPayload("book", "ETHBTC"),
		onBitfinexBookMessage("ethbtc"))
	
	monitorWebsocket(
		bitfinexWS,
		getPayload("trades", "tBTCUSD"),
		onBitfinexTradeMessage)
}

func getPayload(channel string, pair string) string {
	return 	`{"event": "subscribe", "channel": "` + channel + `", "pair": "` + pair + `"}`
}

// bitfinex trade stream format:
// "[channel_id_int, event_string, [id_int, milli_time_int, volume_float, price_float]]"
// positive vol means 'buy', neg vol means 'sell'
// (meaningless to specify buy/sell for a trade, prob an artifact from the orderbook)
// there are two events for a trade, first a low-latency event_string='te'
// then a confirmation event_string='tu'
// trade channel is "25"
func onBitfinexTradeMessage(entry string) {
	entry = strings.Replace(entry, "[", "", -1)
	entry = strings.Replace(entry, "]", "", -1)
	parts := strings.Split(entry, ",")

	if len(parts) == 2 && parts[1] == `"hb"` {
		// TODO: heartbeat
	} else if len(parts) == 6 && parts[1] == `"te"` {
		price, err := strconv.ParseFloat(parts[5], 64)
		unixtime, err := strconv.ParseInt(parts[3], 10, 64)
		timestamp := time.Unix(unixtime/1000, unixtime%1000)
		volume, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
		if price > 0 {
			writeBitfinexTradeEntry(price, timestamp, volume)
		}
	} else if len(parts) == 6 && parts[1] == `"tu"` {
		//TODO: order confirmation 	
	} else if len(parts) > 25 {
		//TODO: snapshot
	} else {
		log.Printf("dont understand (%s): %s", len(parts), entry)
	}
}

func writeBitfinexTradeEntry(price float64, timestamp time.Time, volume float64) {
	//log.Printf("price: %s, time: %s, vol: %s", price, timestamp, volume)
	queryStr := "INSERT INTO bitfinex_trades_btcusd(price, volume, time_stamp, time_recieved) VALUES($1, $2, $3, CURRENT_TIMESTAMP);"
	_, err := db.Exec(queryStr, price, volume, timestamp)
	if err != nil {
		// we are inserting a trade that already exists (same timestamp)
		return
	}
	rd.Publish("bitfinex", "trade")
}

type BitfinexBookEntry struct {
	Price float64
	OrderCount int64
	Volume float64
}

// bitfinex book stream format:
// "[channel_id_int,[price_float,count_int,volume_float]]"
func onBitfinexBookMessage(currency string) func(string) {
	return func(entry string) {
		entry = strings.Replace(entry, "[", "", -1)
		entry = strings.Replace(entry, "]", "", -1)
		parts := strings.Split(entry, ",")

		if len(parts) == 2 && parts[1] == `"hb"` {
			// TODO: heartbeat
		} else if len(parts) == 4 {
			price, err := strconv.ParseFloat(parts[1], 64)
			orderCount, err := strconv.ParseInt(parts[2], 10, 64)
			volume, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			if price > 0 {
				writeBitfinexBookEntry(price, orderCount, volume, currency)
			}
		} else if len(parts) > 25 {
			// asks := make([]BitfinexBookEntry, 0, len(parts))
			// bids := make([]BitfinexBookEntry, 0, len(parts))
			rd.Del(currency + "-bitfinex-asks")
			rd.Del(currency + "-bitfinex-bids")
			for i := 1; i < len(parts); i+=3 {
				price, err := strconv.ParseFloat(parts[i], 64)
				orderCount, err := strconv.ParseInt(parts[i+1], 10, 64)
				volume, err := strconv.ParseFloat(parts[i+2], 64)
				if err != nil {
					log.Fatal(err)
				}
				// log.Printf("price: %s, count: %s, vol: %s", price, orderCount, volume)
				writeBitfinexBookEntry(price, orderCount, volume, currency)
				// if volume < 0 {
				// 	asks = append(asks, BitfinexBookEntry{
				// 		Price: price,
				// 		OrderCount: orderCount,
				// 		Volume: volume,
				// 	})
				// } else {
				// 	bids = append(bids, BitfinexBookEntry{
				// 		Price: price,
				// 		OrderCount: orderCount,
				// 		Volume: volume,
				// 	})
				// }
			}
			// resetBitfinexBook(asks, bids, currency)
		} else {
			log.Printf("dont understand (%s): %s", len(parts), entry)
		}
	}
}

func writeBitfinexBookEntry(price float64, orderCount int64, volume float64, currency string) {
	orderType := "buy"
	if volume < 0 {
		// this is an 'ask' order
		orderType = "sell"
		volume *= -1
	}
	queryStr := "INSERT INTO bitfinex_book_" + currency + "(price, order_count, volume, order_type, time_stamp) VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)"
	_, err := db.Exec(queryStr, price, orderCount, volume, orderType)
	if err != nil {
		log.Fatal("insert err", err)
	}

	key := currency + "-bitfinex-bids"
	if orderType == "sell" {
		key = currency + "-bitfinex-asks"
	}

	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	vals, err := rd.ZRangeByScore(key, redis.ZRangeBy{
		Min: priceStr,
		Max: priceStr,
	}).Result()
	if err != nil || len(vals) > 1 {
		log.Fatal("redis err: ", err)
	}
	if len(vals) == 1 {
		entryStr := vals[0]
		rd.ZRem(key, entryStr)
	} 
	
	if orderCount != 0 {
		entry := BitfinexBookEntry{
			Price: price,
			OrderCount: orderCount,
			Volume: volume,
		}
		entryStr, err := json.Marshal(entry)
		if err != nil {
			log.Fatal("json parse error")
		}
		// log.Printf("adding: %s, price: %s, str: %s", key, price, entryStr)
		
		rd.ZAdd(key, redis.Z{
			Score: price,
			Member: entryStr,
		})
	}
	rd.Publish("bitfinex", currency + "book")
}

// func resetBitfinexBook(asks []BitfinexBookEntry, bids []BitfinexBookEntry, currency string) {
// 	rd.Del(currency + "-bitfinex-asks")
// 	for _, ask := range asks {
// 		entryStr, err := json.Marshal(ask)
// 		if err != nil {
// 			log.Fatal("json parse error")
// 		}
	
// 		rd.ZAdd(currency + "-bitfinex-asks", redis.Z{
// 			Score: ask.Price,
// 			Member: entryStr,
// 		})
// 	}
	
// 	rd.Del(currency + "-bitfinex-bids")
// 	for _, bid := range bids {
// 		entryStr, err := json.Marshal(bid)
// 		if err != nil {
// 			log.Fatal("json parse error")
// 		}
	
// 		rd.ZAdd(currency + "-bitfinex-bids", redis.Z{
// 			Score: bid.Price,
// 			Member: entryStr,
// 		})
// 	}
// }
