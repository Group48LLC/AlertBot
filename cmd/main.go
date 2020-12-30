package main

import(
	"fmt"
	// s3 "github.com/Group48LLC/AlertBot/pkg/s3"
	myBinance "github.com/Group48LLC/AlertBot/pkg/myBinance"
	secret "github.com/Group48LLC/AlertBot/pkg/secret"
	
)

func main(){
	user := "test1"

	secretVars, err := secret.GetSecret()
	if err != nil {
		fmt.Println(err)
	}

	// load in specific user keys
	apiKey := secretVars[user + "_api_key"].(string)
	apiSecret := secretVars[user + "_api_secret"].(string)


	// fire off binance web sockets
	listenKey := myBinance.GetListenKey(apiKey, apiSecret)
	myBinance.OpenSocket(listenKey)
	fmt.Println("Main Complete")
}