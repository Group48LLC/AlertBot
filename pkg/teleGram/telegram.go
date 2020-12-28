package telegram

import (
"bytes"
"fmt"
"net/http"
"encoding/json"
)

func sendAlert(text string, botApi string, chatId string) {

	requestUrl := "https://api.telegram.org/" + botApi + "/sendMessage"

	client := &http.Client{}

	values := map[string]string{"text": text, "chatId": chatId }
	jsonParams, _ := json.Marshal(values)

	req, _:= http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if(err != nil){
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
		// defer res.Body.Close()
	}

}

