package myBinance 

import(
	"strconv"
	"fmt"

	"github.com/adshao/go-binance"

	"github.com/Group48LLC/AlertBot/pkg/s3"

)

func GetAccountInfo(apiKey string, secretKey string) []Listing{
	// move to binance pkg //
	// returns a map data structure with assets mapping to three values free locked total
	// string{string float, string float, string float} string{string float, string float, string float}
	returnValue := []Listing{}
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

			returnValue = append(returnValue, Listing{
				Symbol: v.Asset,
				Locked: v.Locked,
				Free: v.Free,
				Total: totalStr,
			})
		}
	}
	return returnValue
}