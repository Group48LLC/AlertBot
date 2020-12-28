package s3

import (
        "fmt"
		"encoding/json"
		"time"
		"bytes"
		"net/http"
		"log"
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/s3"

		secret "github.com/Group48LLC/AlertBot/pkg/secret"

		myBinance "github.com/Group48LLC/AlertBot/pkg/myBinance"

		models "github.com/Group48LLC/AlertBot/pkg/models"
)


func CreateUserUserBalanceData(inputUserId string, data models.UserBalances) models.UserBalanceData{
	returnValue := models.UserBalanceData{
		UserId: inputUserId,
		Balances: data,
	}
	return returnValue
}


func UploadToS3(filePath string, userId string, jsonData []byte) error {
	// uploads account balence data to s3
	// returns nothing or an error if error

	// build these vars from env
	const (
		S3_REGION = "us-east-1"
		S3_BUCKET = "g48-alert-bot-48"
	)

	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
    if err != nil {
        log.Fatal(err)
    }

    // Get file size and read the file content into a buffer
    // Config settings: this is where you choose the bucket, filename, content-type etc.
    // of the file you're uploading.
    _, err = s3.New(s).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET),
        Key:                  aws.String(filePath),
        ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(jsonData),
        ContentType:          aws.String(http.DetectContentType(jsonData)),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"),
	})
    return err
}


func CreateJSONFile(data models.UserBalanceData, filePath string) []byte{
	file, _ := json.MarshalIndent(data, "", " ")
	return file
}


func WriteUserBalanceData(userId string) (string, error) {
	// handles the request to api

	start1 := time.Now()
	// timing aws secret call

	secretVars, err := secret.GetSecret()
	if err != nil {
		return "ERROR", err
	}
	
	
	apiKey := secretVars["test1_api_key"].(string)
	apiSecret := secretVars["test1_api_secret"].(string)

	t1 := time.Now() // end time
	elapsed1 := t1.Sub(start1)
	
	start2 := time.Now()
	// timing binance call

	accountBalances := myBinance.GetAccountInfo(apiKey, apiSecret)
	userIdBalances := CreateUserUserBalanceData(userId, accountBalances)

	t2 := time.Now() // end time
	elapsed2 := t2.Sub(start2)


	start3 := time.Now() // start of s3 write timing

	filePath := "./"
	filePath += userId
	filePath += ".json"
	s3Prefix := userId
	s3Prefix += "Balances"
	jsonData := CreateJSONFile(userIdBalances, filePath)

	UploadToS3(s3Prefix, userId, jsonData)

	t3 := time.Now() // s3 writing end time
	elapsed3 := t3.Sub(start3)
	fmt.Println("\nS3 Writing Complete: ")
	
	fmt.Println(userId)
	fmt.Println("AWS get secret time: ", elapsed1)
	fmt.Println("Binance api time: ", elapsed2)
	fmt.Println("S3 Writing time: ", elapsed3)
	returnValue := "----success----"
	
	return returnValue, nil
}

