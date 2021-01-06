package main
// TODO:
// do first !! Build user model adjust code.
	// track orders T | F track trades T | F IsPremium T | F isAuthed(paid up) T | F 
	// Figure out how customers will get binance keys (manual tutorial?)
		// EMPHASIZE API IS READ ONLY!!
		// create way to detect if api is more than read only? then alert user.
	// create automated way to add to secret (goSDK aws)
	// created automated way to launch telegram bot for customer
	// save all api keys in secret as userid_api_key/secret
	// decide on how im going to split secrets. one for each exchange?
	// Feed in list of user id's
	// setup api keys/vars for each user id
	// open websocket for-each userid
	// create a way to check which user id's do NOT have an open websocket for their exechange(s)
		// add in a way to flag multiple failures and/or test api key for validity or expired
	// create DB to store user data
	// decide what user data to store
	// last but not least create web front-end

import(
	"fmt"
	"time"
	"strconv"
	// s3 "github.com/Group48LLC/AlertBot/pkg/s3"
	myBinance "github.com/Group48LLC/AlertBot/pkg/myBinance"
	secret "github.com/Group48LLC/AlertBot/pkg/secret"
	
)

func monitor(user string, trackOrders bool){
	secretVars, err := secret.GetSecret()
	if err != nil {
		fmt.Println(err)
	}
	apiKey := secretVars[user + "_api_key"].(string)
	apiSecret := secretVars[user + "_api_secret"].(string)
	


	// fire off binance web sockets
	listenKey := myBinance.GetListenKey(apiKey, apiSecret)
	myBinance.OpenSocket(listenKey, trackOrders)
	fmt.Println("Main Complete")
}

func main(){
	users := []string {"test2"}
	userTrackOrders := true
	
	for _, user := range users{
		go monitor(user, userTrackOrders)
	}

	minuteCount := 0
	for true {
		time.Sleep(60 * time.Second)
		minuteCount ++
		fmt.Println("Minutes monitored: " + strconv.Itoa(minuteCount))
	}
}