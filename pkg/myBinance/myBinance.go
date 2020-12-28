package myBinance 

import(
	"strconv"
	"fmt"
	"context"
	// // "io/ioutil"
	"encoding/json"
    // "net/http"
	// "strings"
	"log"
	"github.com/adshao/go-binance"

	models "github.com/Group48LLC/AlertBot/pkg/models"
	telegram "github.com/Group48LLC/AlertBot/pkg/telegram"
	secret "github.com/Group48LLC/AlertBot/pkg/secret"

)



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


func GetListenKey(apiKey string, apiSecret string) string{
	client := binance.NewClient(apiKey, apiSecret)
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
	secretVars, err := secret.GetSecret()
	telegramApi := secretVars["test1_telegram_api"]


	marketCoinQty, err := strconv.ParseFloat(tradeQty, 64)
	// checkError("QTY error", err)
	marketCoinPrice, err := strconv.ParseFloat(price, 64)
	
	marketCoinTotal := marketCoinQty * marketCoinPrice // computes total
	
	checkError("QTY error", err)
	fmt.Println(timeStr)
	fmt.Println("\n ", tradeType, " : ", tradeSymbol, " : ", tradeDetail) // trade type Buy/sell
	fmt.Println("Order Qty", orderQty) // order qty
	fmt.Println("Trade Qty: ", tradeQty) // trade qty
	fmt.Println("Price: ", price) // price
	if (tradeType == "SELL"){
		fmt.Println("You gained ", marketCoinTotal)
		telegram.SendAlert("hello telegram", "botApi<TOKEN>", "<chatId>")
	}
	if (tradeType == "BUY"){
		fmt.Println("You spent ", marketCoinTotal)
	}
}


func OpenSocket(listenKey string){
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
				fmt.Println("Order Created")
			}
			
		}
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsUserDataServe(listenKey, wsHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}