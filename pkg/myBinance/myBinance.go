package myBinance 

import(
	"strconv"
	"context"
	"fmt"
	// // "io/ioutil"
	"encoding/json"
    // "net/http"
	// "strings"
	"log"
	"time"
	binance "github.com/Group48LLC/AlertBot/pkg/goBinance/"
	"github.com/gorilla/websocket"

	models "github.com/Group48LLC/AlertBot/pkg/models"
	telegram "github.com/Group48LLC/AlertBot/pkg/telegram"
	secret "github.com/Group48LLC/AlertBot/pkg/secret"

)

var (
	baseURL         = "wss://stream.binance.us:9443/ws"
	baseFutureURL   = "wss://fstream.binance.us/ws"
	combinedBaseURL = "wss://stream.binance.us:9443/stream?streams="
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
)

type WsConfig struct {
	Endpoint string
}

func GetAccountInfo(apiKey string, secretKey string) models.UserBalances{
	// move to binance pkg //
	// returns a map data structure with assets mapping to three values free locked total
	// string{string float, string float, string float} string{string float, string float, string float}
	returnValue := models.UserBalances{}
	client := binance.NewClient(apiKey, secretKey)
	client.NewSetServerTimeService().Do(context.Background())
	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res.Balances {
		lockedBalance, err := strconv.ParseFloat(v.Locked, 32)
		freeBalance, err := strconv.ParseFloat(v.Free, 32)
		
		if err != nil{
			fmt.Println(err)
		}
		if lockedBalance > 0 || freeBalance > 0 {
			var total float64 = float64(lockedBalance + freeBalance)

			// create func that makes this dynamic ( removes trailing 0's)
			totalStr := strconv.FormatFloat(total, 'f', 8, 64)

			returnValue.Balances = append(returnValue.Balances, models.Balance{
				Symbol: v.Asset,
				Locked: v.Locked,
				Free: v.Free,
				Total: totalStr,
			})
		}
	}
	return returnValue
}


func getRecentTrades(sym string){
	// 
	// Pull recent trades for symbols vs Market coin btc/eth/usdt
	// return recent trade data to main

}


func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

var wsServe = func(cfg *WsConfig, handler binance.WsHandler, errHandler binance.ErrHandler) (doneC, stopC chan struct{}, err error) {
	c, _, err := websocket.DefaultDialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if binance.WebsocketKeepalive {
			keepAlive(c, binance.WebsocketTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func WsUserDataServe(listenKey string, handler binance.WsHandler, errHandler binance.ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", baseURL, listenKey)
	cfg := newWsConfig(endpoint)
	return wsServe(cfg, handler, errHandler)
}


func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}


func GetListenKey(apiKey string, apiSecret string) string{
	// return error also pass error along if thrown
	
	client := binance.NewClient(apiKey, apiSecret)
	// client.BaseURL = "https://api.binance.us"
	fmt.Println(client)
	res, err := client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(res)
	return res
}


func ComputeTotalSpent(a float64, b float64) (float64){
	return a * b 
}


func checkError(message string, err error) {
	// Checks for an error, raises error if error
	// no return value
    if err != nil {
		log.Fatal(message, err)
		// Exit(1)
    }
}


func HandleTrade(result map[string]interface{}){
	// pull off market coin (vs coin) ex eth/btc market coin would be btc
	// pull off asset coin (the coin you are buying/selling) ex eth/btc asset coin would be eth
	// list asset and market coins in output
	time := result["E"].(float64)
	timeStr := strconv.FormatFloat(time, 'f', 6, 64)
	var price string = result["p"].(string)
	var tradeQty string = result["z"].(string)
	var orderQty string = result["q"].(string)
	var tradeType string = result["S"].(string)
	var tradeDetail string = result["o"].(string)
	var tradeSymbol string = result["s"].(string)
	// VAR CREATION move this to config model? or file? look into better way of doing this.
	secretVars, err := secret.GetSecret() // move to somethign different own secret
	telegramApi := secretVars["test1_telegram_api"].(string)


	// marketCoinQty, err := strconv.ParseFloat(tradeQty, 64)
	// // checkError("QTY error", err)
	// marketCoinPrice, err := strconv.ParseFloat(price, 64)
	
	// marketCoinTotal := marketCoinQty * marketCoinPrice // computes total
	
	checkError("QTY error", err)
	fmt.Println(timeStr)
	fmt.Println("\n ", tradeType, " : ", tradeSymbol, " : ", tradeDetail) // trade type Buy/sell
	fmt.Println("Order Qty", orderQty) // order qty
	fmt.Println("Trade Qty: ", tradeQty) // trade qty
	fmt.Println("Price: ", price) // price

	var alertMessage string = `Trade Executed:` + "\n"
	alertMessage += tradeType+ " : " + tradeSymbol + " : " + tradeDetail + "\n"
	alertMessage += "Trade Qty: " + tradeQty + "\n"
	alertMessage += "Order Qty: " + orderQty + "\n"
	alertMessage += "Price: " + price

	chatId := "-1001219499639"
	
	telegram.SendAlert(alertMessage, telegramApi, chatId)

}


func HandleOrder(result map[string]interface{}){
	time := result["E"].(float64)
	timeStr := strconv.FormatFloat(time, 'f', 6, 64)
	var price string = result["p"].(string)
	var orderQty string = result["q"].(string)
	var orderType string = result["S"].(string)
	var orderDetail string = result["o"].(string)
	var orderSymbol string = result["s"].(string)
	var orderAction string = result["X"].(string) // CANCELED | FILLED etc..

	// VAR CREATION move this to config model? or file? look into better way of doing this.
	secretVars, err := secret.GetSecret() // move to somethign different own secret
	telegramApi := secretVars["test1_telegram_api"].(string)
	
	checkError("QTY error", err)
	fmt.Println(timeStr)
	fmt.Println("\n ", orderType, " : ", orderSymbol, " : ", orderDetail) // trade type Buy/sell
	fmt.Println("Order Qty", orderQty) // order qty
	fmt.Println("Price: ", price) // price

	var alertMessage string = ``+ orderAction + " Order: " + "\n"
	alertMessage += orderType+ " : " + orderSymbol + " : " + orderDetail + "\n"
	alertMessage += "Order Qty: " + orderQty + "\n"
	alertMessage += "Price: " + price

	chatId := "-1001219499639"
	
	telegram.SendAlert(alertMessage, telegramApi, chatId)
}


func OpenSocket(listenKey string, trackOrders bool){
	var result map[string]interface{}
	wsHandler := func(message []byte) {
		// fmt.Println(string(message))
		json.Unmarshal(message, &result)
		if result["e"] == "executionReport" { // valid results that ended with a market action (trade or order creation)	
			fillQty, err := strconv.ParseFloat(result["z"].(string), 64) // converts string to float
			if err != nil{
				log.Fatal(err)
			}
			fmt.Println(result)
			if fillQty > 0 { // results that are a valid trade
				HandleTrade(result)
			} else { // Valid order results
				if trackOrders == true {
					HandleOrder(result)
					fmt.Println("Order Created")
				} else {
					fmt.Println("Track Orders: ")
					fmt.Println(trackOrders)
				}
			}
			
		}
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}
	
	doneC, _, err := WsUserDataServe(listenKey, wsHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}