package main

import(
	"fmt"
	// s3 "github.com/Group48LLC/AlertBot/pkg/s3"
	myBinance "github.com/Group48LLC/AlertBot/pkg/myBinance"
	secret "github.com/Group48LLC/AlertBot/pkg/secret"
	
)

func main(){
	secretVars, err := secret.GetSecret()
	if err != nil {
		fmt.Println(err)
	}
	apiKey := secretVars["test1_api_key"].(string)
	apiSecret := secretVars["test1_api_secret"].(string)


	// fire off binance web sockets
	listenKey := myBinance.GetListenKey(apiKey, apiSecret)
	myBinance.OpenSocket(listenKey)
	fmt.Println("Main Complete")
}